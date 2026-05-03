package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"
)

var MODE string

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
			for _, allowed := range allowedOrigins {
				allowed = strings.TrimSpace(allowed)
				if allowed != "" && origin == allowed {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					break
				}
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
	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
		log.Fatal("Error loading .env file")
		return
	}
	MODE = os.Getenv("MODE")
	if MODE == "local" {
		os.Setenv("host", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
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
	router.HandleFunc("/settings/retention", changeRetention)
	router.HandleFunc("/service/congestion", getCongestion)
	go StartWebSocketServer(router)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
