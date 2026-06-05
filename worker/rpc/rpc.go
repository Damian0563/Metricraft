package rpc

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"worker/db"
)

type Server struct {
	pb.UnimplementedMetricraftServer
}

func (s *Server) GetTrafficCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	start := req.Start.AsTime()
	resolution := req.Resolution
	return db.GetTrafficCongestion(ctx, start, resolution)
}
