package rpc

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"worker/db"
	"worker/worker"
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

func (s *Server) GetGeographicalPerformance(ctx context.Context, req *pb.Timeframe) (*pb.FloatDistribution, error) {
	return db.GetGeographicalPerformance(ctx, req.Start.AsTime(), req.Timezone)
}

func (s *Server) GetStatusCodeDistribution(ctx context.Context, req *pb.Timeframe) (*pb.Distribution, error) {
	return nil, nil
}

func (s *Server) UpdateWorker(ctx context.Context, req *pb.Worker) (*pb.Status, error) {
	status := worker.TestWorker(ctx, req)
	worker.RegisterAndStartWorker(req)
	return status, nil
}

func (s *Server) DeleteWorker(ctx context.Context, req *pb.WorkerUrl) (*pb.Status, error) {
	if err := db.DeleteWorkerlogs(ctx, req.Url); err != nil {
		return nil, err
	}
	worker.CancelWorker(req.Url)
	return &pb.Status{Success: true}, nil
}

func (s *Server) CreateWorker(ctx context.Context, req *pb.Worker) (*pb.Status, error) {
	status := worker.TestWorker(ctx, req)
	worker.RegisterAndStartWorker(req)
	return status, nil
}

func (s *Server) GetWorkerUptime(ctx context.Context, req *pb.WorkerUrl) (*pb.WorkerUptime, error) {
	return db.GetWorkerUptime(ctx, req.Url, req.Timezone, req.PollInterval)
}
