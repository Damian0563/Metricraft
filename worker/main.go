package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	pb "metricraft/proto/metricraft/proto"
	"metricraft/worker/enter"
	"net"
	"net/http"
	"os"
)

func main() {
	var err error
	err = godotenv.Load()
	MODE := os.Getenv("MODE")
	if MODE == "local" {
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
	router := http.NewServeMux()
	router.HandleFunc("/", enter.Enter)
	ctx := context.Background()
	errChannel := make(chan error)
	go enter.InitDB(ctx, errChannel)
	err = <-errChannel
	if err != nil {
		panic(err)
	}
	go func() {
		lis, err := net.Listen("tcp", "localhost:50051")
		if err != nil {
			panic(err)
		}
		grpcServer := grpc.NewServer()
		s := &server{}
		s.loadFeatures()
		pb.RegisterMetricraftServer(grpcServer, s)
		fmt.Println("gRPC server listening on port 50051")
		grpcServer.Serve(lis)
	}()
	fmt.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", router); err != nil {
		panic(err)
	}
}
