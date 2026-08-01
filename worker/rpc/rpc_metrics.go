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
	return db.GetTrafficCongestion(ctx, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetGeographicalTraffic(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetGeographicalTraffic(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetP95Latency(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetP95Latency(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetUptimeScore(ctx context.Context, req *pb.Timeframe) (*pb.FloatDistribution, error) {
	return db.GetUptimeScore(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetThroughput(ctx context.Context, req *pb.Timeframe) (*pb.Throughput, error) {
	return db.GetThroughput(ctx, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetGeographicalPerformance(ctx context.Context, req *pb.Timeframe) (*pb.FloatDistribution, error) {
	return db.GetGeographicalPerformance(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetStatusCodeDistribution(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetStatusCodeDistribution(ctx, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetRouteCongestion(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return db.GetRouteCongestion(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetHttpMethodDistribution(ctx context.Context, req *pb.Timeframe) (*pb.Congestion, error) {
	return db.GetHttpMethodDistribution(ctx, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetUniqueVisitors(ctx context.Context, req *pb.Timeframe) (*pb.SimpleRepeatedDistribution, error) {
	return db.GetUniqueVisitors(ctx, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}

func (s *Server) GetHotHours(ctx context.Context, req *pb.Timeframe) (*pb.SimpleRepeatedDistribution, error) {
	return db.GetHotHours(ctx, req.Rules, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetCustomMetricData(ctx context.Context, req *pb.CustomMetricRequest) (*pb.CustomMetricData, error) {
	return db.GetCustomMetricData(ctx, req.Metric, req.Rules, req.Start.AsTime(), req.Resolution, req.Timezone)
}
