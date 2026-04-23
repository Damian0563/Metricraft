package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func getConfig() map[string]string {
	response := make(map[string]string)
	body, _ := os.ReadFile("../config.json")
	json.Unmarshal(body, &response)
	return response
}

var config map[string]string = getConfig()
var PORT string = config["port"]
var MODE string = config["mode"]

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		host := os.Getenv("host")
		if host != "" {
			ws_host := os.Getenv("ws")
			if origin == host+":"+PORT || origin == host+":8081" || origin == ws_host+":8080/ws/visitors" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		} else {
			frontend := os.Getenv("frontend")
			worker := os.Getenv("worker")
			ws := os.Getenv("ws")
			if origin == frontend+":"+PORT || origin == frontend+":8081" || origin == worker+":8081" || origin == ws+":8080/ws/visitors" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, session-token")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	var err error
	if MODE == "local" {
		os.Setenv("host", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		err = godotenv.Load("../.env")
		if err != nil {
			fmt.Println(err)
			log.Fatal("Error loading .env file")
			return
		}
	} else if MODE == "docker" {
		os.Setenv("frontend", "http://metricraft-frontend-1")
		os.Setenv("worker", "http://metricraft-metricraft-1")
		os.Setenv("ws", "ws://metrcraft-backend-1")
		os.Setenv("redis", "metricraft-redis-1")
	} else {
		os.Setenv("frontend", "https://metricraft-metricraft-1")
		os.Setenv("worker", "https://metricraft-worker-1")
		os.Setenv("ws", "wss://metrcraft-backend-1")
		os.Setenv("redis", "metricraft-redis-1")
	}
	router := http.NewServeMux()
	router.HandleFunc("/", welcome)
	router.HandleFunc("/sign", sign)
	router.HandleFunc("/dashboard/init", dashboardInit)
	router.HandleFunc("/settings/realtime", toggleRealtime)
	go StartWebSocketServer(router)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
