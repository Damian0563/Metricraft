package main

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"metricraft/worker/enter"
)

type server struct {
	pb.UnimplementedMetricraftServer
}

func (s *server) GetTrafficCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	start := req.Start.AsTime()
	resolution := req.Resolution
	return enter.GetTrafficCongestion(ctx, start, resolution)
}
