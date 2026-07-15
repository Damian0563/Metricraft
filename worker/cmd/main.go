package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	pb "metricraft/proto/metricraft/proto"
	"net"
	"net/http"
	"os"
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

func main() {
	loadEnv()
	MODE := os.Getenv("MODE")
	if os.Getenv("SECRET") == "" || os.Getenv("DATABASE_LOGS") == "" {
		panic("Secret must be set")
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
