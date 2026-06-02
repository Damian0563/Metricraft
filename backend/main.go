package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"time"
)

var MODE string
var conn *grpc.ClientConn

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
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
	godotenv.Load()
	MODE = os.Getenv("MODE")
	if MODE == "local" {
		os.Setenv("host", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		os.Setenv("grpc", "127.0.0.1:50051")
		os.Setenv("frontend", "http://localhost")
		os.Setenv("worker", "http://localhost")
	} else if MODE == "docker" {
		os.Setenv("frontend", "http://metricraft")
		os.Setenv("worker", "http://worker")
		os.Setenv("ws", "ws://backend")
		os.Setenv("redis", "redis:6379")
		os.Setenv("grpc", "worker:50051")
	} else {
		os.Setenv("frontend", "http://localhost")
		os.Setenv("worker", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		os.Setenv("grpc", "127.0.0.1:50051")
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
