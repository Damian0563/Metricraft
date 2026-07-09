package worker

import (
	supabase "backend/db"
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

func TestWorker(ctx context.Context, worker *pb.Worker) *pb.Status {
	req, err := http.NewRequest("GET", worker.Url, nil)
	if err != nil {
		return &pb.Status{Success: false, Err: err.Error(), StatusCode: 500}
	}
	for k, v := range worker.Headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: time.Second * 10}
	resp, err := client.Do(req)
	if err != nil {
		return &pb.Status{Success: false, Err: err.Error(), StatusCode: 500}
	}
	defer resp.Body.Close()
	fmt.Println(worker.Url, resp.StatusCode)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &pb.Status{Success: true, Err: "", StatusCode: int32(resp.StatusCode)}
	} else {
		return &pb.Status{Success: false, Err: fmt.Sprintf("Faulty response code of the health endpoint: %d. Confirm that the endpoint is working and the URL is correct.", resp.StatusCode), StatusCode: int32(resp.StatusCode)}
	}
}

func StartWorker(ctx context.Context, worker *pb.Worker, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
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
		case <-time.After(interval):
			resp := TestWorker(ctx, worker)
			if err := db.InsertWorkerLog(ctx, worker.Url, resp.Success); err != nil {
				fmt.Println(err, worker.Url)
				return
			}
		}
	}
}

func registerAndStart(worker *pb.Worker) {
	CancelWorker(worker.Url)
	Orchestrator.Mutex.Lock()
	workerCtx, cancel := context.WithCancel(context.Background())
	Orchestrator.Registry[worker.Url] = cancel
	Orchestrator.Mutex.Unlock()
	go StartWorker(workerCtx, worker, nil)
}

func RegisterAndStartWorker(worker *pb.Worker) {
	registerAndStart(worker)
}

func CancelWorker(url string) {
	Orchestrator.Mutex.Lock()
	if cancel, ok := Orchestrator.Registry[url]; ok {
		cancel()
		delete(Orchestrator.Registry, url)
	}
	Orchestrator.Mutex.Unlock()
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
	for _, w := range workers {
		workerCtx, cancel := context.WithCancel(ctx)
		Orchestrator.Mutex.Lock()
		if prev, ok := Orchestrator.Registry[w.Url]; ok {
			prev()
		}
		Orchestrator.Registry[w.Url] = cancel
		Orchestrator.Mutex.Unlock()
		wg.Add(1)
		go StartWorker(workerCtx, &pb.Worker{Url: w.Url, PollInterval: int32(w.PollInterval), Headers: w.Headers}, &wg)
	}
	wg.Wait()
}
