package rpc

import (
	"context"
	pb "metricraft/proto/metricraft/proto"
	"worker/db"
	"worker/worker"
)

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
