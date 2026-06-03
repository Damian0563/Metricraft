package rpc

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"worker/db"
	"worker/types"
)

func (s *types.Server) GetTrafficCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	start := req.Start.AsTime()
	resolution := req.Resolution
	return db.GetTrafficCongestion(ctx, start, resolution)
}
