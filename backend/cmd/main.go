package main

import (
	"backend/api"
	"backend/db"
	"backend/redis"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"os"
	"time"
)

var MODE string
var conn *grpc.ClientConn

func loadEnv() {
	for _, path := range []string{".env", "../.env", "backend/.env"} {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}
}

func initPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}
	config.MinConns = 3
	config.MaxConns = 20
	config.MaxConnLifetime = 30 * time.Minute
	config.HealthCheckPeriod = 5 * time.Minute
	config.ConnConfig.ConnectTimeout = 15 * time.Second
	return pgxpool.NewWithConfig(ctx, config)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
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
	loadEnv()
	MODE = os.Getenv("MODE")
	if MODE == "local" {
		os.Setenv("host", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		os.Setenv("grpc", "127.0.0.1:50051")
		os.Setenv("frontend", "http://localhost")
		os.Setenv("worker", "http://localhost")
	} else {
		os.Setenv("frontend", "http://localhost")
		os.Setenv("worker", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("redis", "127.0.0.1:6379")
		os.Setenv("grpc", "127.0.0.1:50051")
	}
	redis.InitClient()
	usersPool, err := initPool(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		fmt.Println("failed to initialize users database pool:", err)
		return
	}
	defer usersPool.Close()
	logsPool, err := initPool(context.Background(), os.Getenv("DATABASE_LOGS"))
	if err != nil {
		fmt.Println("failed to initialize logs database pool:", err)
		return
	}
	defer logsPool.Close()
	db.SetUsersPool(usersPool)
	db.SetLogsPool(logsPool)
	router := chi.NewRouter()
	router.Use(corsMiddleware)
	router.Use(httprate.LimitByIP(20, 1*time.Second))
	router.Get("/welcome", api.Welcome)
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/sign", api.Sign)
	})
	router.Get("/dashboard/init", api.DashboardInit)
	router.Post("/settings/realtime", api.ToggleRealtime)
	router.Post("/settings/retention", api.ChangeRetention)
	router.Post("/settings/metrics", api.ChangeMetricsHandler)
	router.Get("/dashboard/fetch", api.Navigator) //?persist=boolean
	router.Get("/dashboard/worker/list", api.ListWorkers)
	router.Post("/dashboard/worker/new", api.SaveWorker)
	router.Patch("/dashboard/worker/update", api.UpdateWorker)
	router.Delete("/dashboard/worker/delete", api.DeleteWorker)
	router.Post("/dashboard/worker/uptime", api.GetWorkerUptime)
	router.Patch("/dashboard/worker/notifications", api.SaveNotificationRecipients)
	router.Get("/invites/pending", api.PendingInvites)
	router.Get("/invites/team", api.TeamMembers)
	router.Post("/invites/handle", api.HandleInvite) //?user=string&action=boolean
	router.Post("/invites/manual", api.SendInvites)
	router.Post("/invites/batch", api.UploadUsersFromCSV)
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(6, time.Minute))
		r.Post("/verify/send", api.SendVerification)
		r.Post("/verify/check", api.CheckVerification)
	})
	router.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(4, time.Minute))
		r.Post("/recovery/send", api.SendRecovery)
		r.Post("/recovery/check", api.CheckRecovery) //?id=string
	})
	go api.StartWebSocketServer(router)
	conn, err = grpc.NewClient(fmt.Sprintf("dns:///%v", os.Getenv("grpc")), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	api.SetGRPCConn(conn)
	defer conn.Close()
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", corsMiddleware(router))
}
