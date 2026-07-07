package worker

import (
	supabase "backend/db"
	btype "backend/types"
	"context"
	"fmt"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"sync"
	"time"
	"worker/db"
	"worker/types"
)

var Orchestrator types.Orchestrator = types.Orchestrator{
	Mutex:    sync.Mutex{},
	Registry: make(map[string]context.CancelFunc),
}

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
	defer func() {
		Orchestrator.Mutex.Lock()
		delete(Orchestrator.Registry, worker.Url)
		Orchestrator.Mutex.Unlock()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		default:
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
}

func UpdateWorker(worker btype.Worker) {
	CancelWorker(worker.Url)
	Orchestrator.Mutex.Lock()
	context, cancel := context.WithCancel(context.Background())
	Orchestrator.Registry[worker.Url] = cancel
	Orchestrator.Mutex.Unlock()
	go StartWorker(context, &pb.Worker{Url: worker.Url, PollInterval: int32(worker.PollInterval), Headers: worker.Headers})
}

func CancelWorker(url string) {
	cancel, ok := Orchestrator.Registry[url]
	if ok {
		cancel()
		Orchestrator.Mutex.Lock()
		delete(Orchestrator.Registry, url)
		Orchestrator.Mutex.Unlock()
	}
}

func OrchestrateWorkers(ctx context.Context) {
	appName, err := db.GetAppname(ctx)
	if err != nil {
		panic(err)
	}
	workers, err := supabase.GetWorkers(appName)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	for _, worker := range workers {
		ctx, cancel := context.WithCancel(ctx)
		Orchestrator.Mutex.Lock()
		CancelWorker(worker.Url)
		Orchestrator.Registry[worker.Url] = cancel
		Orchestrator.Mutex.Unlock()
		wg.Add(1)
		go StartWorker(ctx, &pb.Worker{Url: worker.Url, PollInterval: int32(worker.PollInterval), Headers: worker.Headers})
		defer wg.Done()
	}
	wg.Wait()
	defer func() {
		Orchestrator.Mutex.Lock()
		for _, cancel := range Orchestrator.Registry {
			cancel()
		}
		Orchestrator.Registry = make(map[string]context.CancelFunc)
		Orchestrator.Mutex.Unlock()
	}()
}
