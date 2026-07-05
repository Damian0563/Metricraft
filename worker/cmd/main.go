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
	for _, path := range []string{".env", "../.env", "worker/.env"} {
		if err := godotenv.Load(path); err == nil {
			return
		}
	}
}

func main() {
	var err error
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
	errChannel := make(chan error)
	go db.InitDB(ctx, errChannel)
	go worker.OrchestrateWorkers(ctx, errChannel)
	go func(errChannel chan error) {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			errChannel <- err
			return
		}
		grpcServer := grpc.NewServer()
		s := &rpc.Server{}
		pb.RegisterMetricraftServer(grpcServer, s)
		fmt.Println("gRPC server listening on port 50051")
		err = grpcServer.Serve(lis)
		if err != nil {
			errChannel <- err
		}
	}(errChannel)
	err = <-errChannel
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}
}
