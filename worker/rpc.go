package main

import (
	"context"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
)

type server struct {
	pb.UnimplementedMetricraftServer
	features map[string]func(context.Context, *pb.Timeframe) (*pb.Congestion, error)
}

func (s *server) loadFeatures() {
	s.features = make(map[string]func(context.Context, *pb.Timeframe) (*pb.Congestion, error))
}

func (s *server) GetTrafficCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	fmt.Println(req)
	return &pb.Congestion{}, nil
}
