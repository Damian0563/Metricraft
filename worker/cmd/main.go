package main

import (
	backenddb "backend/db"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	pb "metricraft/proto/metricraft/proto"
	"net"
	"net/http"
	"os"
	"time"
	"worker/db"
	"worker/enter"
	"worker/rpc"
	"worker/worker"
)

func loadEnv() {
	for _, path := range []string{".env", "../.env", "worker/.env", "backend/.env", "../backend/.env"} {
		_ = godotenv.Load(path)
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

func main() {
	loadEnv()
	MODE := os.Getenv("MODE")
	if os.Getenv("SECRET") == "" || os.Getenv("DATABASE_LOGS") == "" || os.Getenv("DATABASE_USERS") == "" {
		panic("SECRET, DATABASE_LOGS, and DATABASE_USERS must be set")
	}
	if MODE == "local" {
		os.Setenv("backend", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("grpc", "127.0.0.1:50051")
	} else {
		os.Setenv("backend", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("grpc", "127.0.0.1:50051")
	}
	router := http.NewServeMux()
	router.HandleFunc("/", enter.Enter)
	ctx := context.Background()
	logsPool, err := initPool(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		panic(err)
	}
	defer logsPool.Close()
	db.SetLogsPool(logsPool)
	usersPool, err := initPool(ctx, os.Getenv("DATABASE_USERS"))
	if err != nil {
		panic(err)
	}
	defer usersPool.Close()
	backenddb.SetUsersPool(usersPool)
	initErr := make(chan error, 1)
	go db.InitDB(ctx, initErr)
	if err := <-initErr; err != nil {
		panic(err)
	}
	go worker.OrchestrateWorkers(ctx)
	go func() {
		lis, err := net.Listen("tcp", os.Getenv("grpc"))
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		s := &rpc.Server{}
		pb.RegisterMetricraftServer(grpcServer, s)
		fmt.Println("gRPC server listening on port 50051")
		if err := grpcServer.Serve(lis); err != nil {
			panic(err)
		}
	}()
	fmt.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}
}
