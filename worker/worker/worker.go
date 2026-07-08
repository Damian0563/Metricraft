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

func StartWorker(ctx context.Context, worker *pb.Worker, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
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
		case <-time.After(interval):
			resp, err := TestWorker(ctx, worker)
			if err != nil {
				fmt.Println(err, worker.Url)
				return
			}
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
	defer func() {
		Orchestrator.Mutex.Lock()
		for _, cancel := range Orchestrator.Registry {
			cancel()
		}
		Orchestrator.Registry = make(map[string]context.CancelFunc)
		Orchestrator.Mutex.Unlock()
	}()
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
