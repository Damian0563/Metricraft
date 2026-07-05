package worker

import (
	"context"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"time"
)

func TestWorker(ctx context.Context, worker *pb.Worker) (*pb.Status, error) {
	req, err := http.NewRequest("GET", worker.Url, nil)
	for k, v := range worker.Headers {
		req.Header.Set(k, v)
	}
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &pb.Status{Success: true, Err: ""}, nil
	} else {
		return &pb.Status{Success: false, Err: fmt.Sprintf("Faulty response code of the health endpoint: %d. Confirm that the endpoint is working and the URL is correct.", resp.StatusCode)}, nil
	}
}

func StartWorker(ctx context.Context, worker *pb.Worker) {
	if _, err := TestWorker(ctx, worker); err != nil {
		return
	}
	interval := time.Duration(worker.PollInterval) * time.Minute
	for {
		res, err := TestWorker(ctx, worker)
		fmt.Println(res, err)
		time.Sleep(interval)

	}
}
