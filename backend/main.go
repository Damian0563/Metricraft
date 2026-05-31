package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var MODE string
var conn *grpc.ClientConn

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
	if err = godotenv.Load(); err != nil {
		fmt.Println("No .env file loaded, relying on environment variables:", err)
	}
	MODE = os.Getenv("MODE")
	if MODE == "local" {
		os.Setenv("host", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		os.Setenv("grpc", "127.0.0.1:50051")
		if err != nil {
			fmt.Println(err)
			log.Fatal("Error loading .env file")
			return
		}
	} else if MODE == "docker" {
		os.Setenv("frontend", "http://metricraft-metricraft-1")
		os.Setenv("worker", "http://metricraft-worker-1")
		os.Setenv("ws", "ws://metricraft-backend-1")
		os.Setenv("redis", "metricraft-redis-1:6379")
		os.Setenv("grpc", "metricraft-worker-1:50051")
	} else {
		os.Setenv("frontend", "https://metricraft-metricraft-1")
		os.Setenv("worker", "https://metricraft-worker-1")
		os.Setenv("ws", "wss://metricraft-backend-1")
		os.Setenv("redis", "metricraft-redis-1:6379")
		os.Setenv("grpc", "metricraft-worker-1:50051")
	}
	router := chi.NewRouter()
	router.Use(corsMiddleware)
	router.Use(httprate.LimitByIP(20, 1*time.Second))
	router.Get("/", welcome)
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/sign", sign)
	})
	router.Get("/dashboard/init", dashboardInit)
	router.Post("/settings/realtime", toggleRealtime)
	router.Post("/settings/retention", changeRetention)
	router.Post("/settings/metrics", changeMetrics)
	router.Get("/dashboard/fetch", navigator)
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(6, time.Minute))
		r.Post("/verify/send", sendVerification)
		r.Post("/verify/check", checkVerification)
	})
	go StartWebSocketServer(router)
	conn, err = grpc.NewClient(fmt.Sprintf("dns:///%v", os.Getenv("grpc")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return
	}
	defer conn.Close()
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
