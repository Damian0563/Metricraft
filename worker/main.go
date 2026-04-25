package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"metricraft/worker/enter"
	"net/http"
	"os"
)

func getConfig() map[string]string {
	response := make(map[string]string)
	body, _ := os.ReadFile("../config.json")
	json.Unmarshal(body, &response)
	return response
}
func main() {
	MODE := getConfig()["mode"]
	var err error
	if MODE == "local" {
		err = godotenv.Load("../.env")
		if err != nil {
			panic(err)
		}
		if os.Getenv("SECRET") == "" || os.Getenv("DATABASE_LOGS") == "" {
			panic("Secret must be set")
		}
		os.Setenv("backend", "http://localhost")
		os.Setenv("ws", "ws://localhost")
	} else if MODE == "docker" {
		os.Setenv("backend", "http://metricraft-backend-1")
		os.Setenv("ws", "ws://metrcraft-backend-1")
	} else {
		os.Setenv("ws", "wss://metrcraft-backend-1")
		os.Setenv("backend", "https://metricraft-metricraft-1")
	}
	os.Setenv("DEST_PORT", getConfig()["dest-port"])
	router := http.NewServeMux()
	router.HandleFunc("/", enter.Enter)
	ctx := context.Background()
	errChannel := make(chan error)
	go enter.InitDB(ctx, errChannel)
	err = <-errChannel
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port 8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
