package worker

import (
	supabase "backend/db"
	"context"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"time"
	"worker/db"
)

func TestWorker(ctx context.Context, worker *pb.Worker) (*pb.Status, error) {
	req, err := http.NewRequest("GET", worker.Url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range worker.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println(worker.Url, resp.StatusCode)
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
		resp, err := TestWorker(ctx, worker)
		if err != nil {
			fmt.Println(err, worker.Url)
			return
		}
		if err := db.InsertWorkerLog(ctx, worker.Url, resp.Success); err != nil {
			fmt.Println(err, worker.Url)
			return
		}
		time.Sleep(interval)

	}
}

func OrchestrateWorkers(ctx context.Context) {
	time.Sleep(15 * time.Second)
	appName, err := db.GetAppname(ctx)
	if err != nil {
		panic(err)
	}
	workers, err := supabase.GetWorkers(appName)
	if err != nil {
		panic(err)
	}
	for _, worker := range workers {
		go StartWorker(ctx, &pb.Worker{Url: worker.Url, PollInterval: int32(worker.PollInterval), Headers: worker.Headers})
	}
}
