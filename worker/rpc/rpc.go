package rpc

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"time"
	"worker/db"
)

type Server struct {
	pb.UnimplementedMetricraftServer
}

func (s *Server) GetTrafficCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	start := req.Start.AsTime().Add(time.Hour).Truncate(time.Hour)
	resolution := req.Resolution
	return db.GetTrafficCongestion(ctx, start, resolution)
}

func (s *Server) GetGeographicalTraffic(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	start := req.Start.AsTime()
	return db.GetGeographicalTraffic(ctx, start)
}

func (s *Server) GetP95Latency(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	start := req.Start.AsTime()
	return db.GetP95Latency(ctx, start)
}

func (s *Server) GetUptimeScore(ctx context.Context, req *pb.Timeframe) (*pb.FloatDistribution, error) {
	start := req.Start.AsTime()
	return db.GetUptimeScore(ctx, start)
}

func (s *Server) GetThroughput(ctx context.Context, req *pb.Timeframe) (*pb.Throughput, error) {
	start := req.Start.AsTime().Add(time.Hour).Truncate(time.Hour)
	resolution := req.Resolution
	return db.GetThroughput(ctx, start, resolution)
}
