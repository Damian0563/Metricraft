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
	return db.GetTrafficCongestion(ctx, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetGeographicalTraffic(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetGeographicalTraffic(ctx, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetP95Latency(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetP95Latency(ctx, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetUptimeScore(ctx context.Context, req *pb.Timeframe) (*pb.FloatDistribution, error) {
	return db.GetUptimeScore(ctx, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetThroughput(ctx context.Context, req *pb.Timeframe) (*pb.Throughput, error) {
	return db.GetThroughput(ctx, req.Start.AsTime(), req.Resolution, req.Timezone)
}
