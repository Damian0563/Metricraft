package main

import (
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
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
	if os.Getenv("SECRET") == "" {
		panic("Secret must be set")
	}
	os.Setenv("DEST_PORT", getConfig()["dest-port"])
	router := http.NewServeMux()
	router.HandleFunc("/", enter.Enter)
	fmt.Println("Listening on port 8081")
	err = http.ListenAndServe(":8081", router)
	if err != nil {
		panic(err)
	}
}
