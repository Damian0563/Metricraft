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
	godotenv.Load()
	MODE := os.Getenv("MODE")
	if os.Getenv("SECRET") == "" || os.Getenv("DATABASE_LOGS") == "" {
		panic("Secret must be set")
	}
	if MODE == "local" {
		os.Setenv("backend", "http://localhost")
		os.Setenv("ws", "ws://localhost")
		os.Setenv("grpc", "127.0.0.1:50051")
	} else if MODE == "docker" {
		os.Setenv("backend", "http://backend")
		os.Setenv("ws", "ws://backend")
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
	go enter.InitDB(ctx, errChannel)
	go func(errChannel chan error) {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			errChannel <- err
			return
		}
		grpcServer := grpc.NewServer()
		s := &server{}
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
