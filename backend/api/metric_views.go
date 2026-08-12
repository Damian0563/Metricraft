package api

import (
	"backend/auth"
	"backend/db"
	"backend/types"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"strings"
	"sync"
)

var grpcConn *grpc.ClientConn

func SetGRPCConn(conn *grpc.ClientConn) {
	grpcConn = conn
}

func Navigator(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	metric := r.URL.Query().Get("metric")
	timeframe := r.URL.Query().Get("timeframe")
	timezone := r.URL.Query().Get("timezone")
	persistChan := make(chan error, 1)
	if persist := r.URL.Query().Get("persist"); persist == "true" {
		go db.PersistTimeframeSelection(persistChan, metric, timeframe)
	} else {
		persistChan <- nil
	}
	convertedTimeframe, resolution := convertTimeframe(timeframe, timezone)
	var response any
	var err error
	ctx := context.Background()
	rules, err := db.GetRules(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var ruleInput []*pb.Rule
	for _, rule := range rules {
		ruleInput = append(ruleInput, &pb.Rule{Rule: rule.Rule, Mode: rule.Mode})
	}
	timeframeEntry := &pb.Timeframe{Rules: ruleInput, Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days, Timezone: timezone}
	switch metric {
	case "Traffic congestion trends":
		response, err = client.GetTrafficCongestion(ctx, timeframeEntry)
	case "Geographical traffic":
		response, err = client.GetGeographicalTraffic(ctx, timeframeEntry)
	case "P95 Latency":
		response, err = client.GetP95Latency(ctx, timeframeEntry)
	case "Uptime Score":
		response, err = client.GetUptimeScore(ctx, timeframeEntry)
	case "Throughput":
		response, err = client.GetThroughput(ctx, timeframeEntry)
	case "Geographic performance":
		response, err = client.GetGeographicalPerformance(ctx, timeframeEntry)
	case "Status code distribution":
		response, err = client.GetStatusCodeDistribution(ctx, timeframeEntry)
	case "Route congestion":
		response, err = client.GetRouteCongestion(ctx, timeframeEntry)
	case "HTTP method distribution":
		response, err = client.GetHttpMethodDistribution(ctx, timeframeEntry)
	case "Unique visitors":
		response, err = client.GetUniqueVisitors(ctx, timeframeEntry)
	case "Hot hours":
		response, err = client.GetHotHours(ctx, timeframeEntry)
	default:
		break
	}
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = <-persistChan //this can be swallowed internally, even if error occured
	httpresponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(httpresponse)
}

func CustomMetricFetch(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	timezone := r.URL.Query().Get("timezone")
	ctx := context.Background()
	metrics, err := db.ListCustomMetrics(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rules, err := db.GetRules(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var ruleInput []*pb.Rule
	for _, rule := range rules {
		ruleInput = append(ruleInput, &pb.Rule{Rule: rule.Rule, Mode: rule.Mode})
	}
	client := pb.NewMetricraftClient(grpcConn)
	var mu sync.Mutex
	var wg sync.WaitGroup
	metricData := make([]types.MetricData, 0, len(metrics))
	var errors []string
	for _, metric := range metrics {
		wg.Add(1)
		go func(metric db.CustomMetric) {
			defer wg.Done()
			standardizedTimeframe := standardizeTimeframe(metric.Timeframe)
			start, resolution := convertTimeframe(standardizedTimeframe, timezone)
			metricRpc := &pb.CustomMetricRequest{
				Metric: &pb.CustomMetric{
					Name:        metric.Name,
					Method:      metric.Method,
					Path:        metric.Path,
					Source:      metric.Source,
					Selector:    metric.Selector,
					Aggregation: metric.Aggregation,
					Timeframe:   metric.Timeframe,
					ValueType:   metric.ValueType,
					ApplyRules:  metric.ApplyRules,
					ChartType:   metric.ChartType,
				},
				Start:      timestamppb.New(start),
				Resolution: resolution.Days,
				Timezone:   timezone,
				Rules:      ruleInput,
			}
			resp, err := client.GetCustomMetricData(ctx, metricRpc)
			mu.Lock()
			if err != nil {
				errors = append(errors, err.Error()+"\n")
			}
			accumulate := false
			if metric.ChartType == "pie" {
				accumulate = true
			}
			metricData = append(metricData, types.MetricData{
				Name:          metric.Name,
				Metrics:       customMetricDataFromProto(resp),
				Timeframe:     metric.Timeframe,
				CustomMetrics: true,
				Accumulate:    accumulate,
			})
			mu.Unlock()
		}(metric)
	}
	wg.Wait()
	w.Header().Set("Content-Type", "application/json")
	response := types.MetricDataResponse{Metrics: metricData, Errors: errors}
	httpresponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(httpresponse)
}

func SaveWorker(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	var worker types.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil || worker.PollInterval < 10 || worker.PollInterval > 60 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	errChan := make(chan error)
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	go db.SaveWorker(appName, worker, errChan)
	response, err := client.CreateWorker(context.Background(), &pb.Worker{
		Url:          worker.Url,
		PollInterval: int32(worker.PollInterval),
		Headers:      worker.Headers,
		AppName:      appName,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := <-errChan; err != nil {
		if err.Error() == "Worker limit reached." || err.Error() == "Worker already exists for this url." {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if httpresponse, err := json.Marshal(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Write(httpresponse)
	}
}

func ListWorkers(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	workers, err := db.GetWorkers(appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpresponse, err := json.Marshal(workers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(httpresponse)
}

func DeleteWorker(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	type dummyWorker struct {
		Url string `json:"url"`
	}
	var dummy dummyWorker
	if err := json.NewDecoder(r.Body).Decode(&dummy); err != nil || strings.TrimSpace(dummy.Url) == "" {
		http.Error(w, "Invalid worker URL", http.StatusBadRequest)
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := db.DeleteWorker(appName, dummy.Url); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	response, err := client.DeleteWorker(context.Background(), &pb.WorkerUrl{Url: dummy.Url})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpresponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(httpresponse)
}

func UpdateWorker(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	var worker types.Worker
	if err := json.NewDecoder(r.Body).Decode(&worker); err != nil || worker.PollInterval < 10 || worker.PollInterval > 60 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := db.UpdateWorker(appName, worker); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	response, err := client.UpdateWorker(context.Background(), &pb.Worker{
		Url:          worker.Url,
		PollInterval: int32(worker.PollInterval),
		Headers:      worker.Headers,
		AppName:      appName,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpresponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(httpresponse)
}

func GetWorkerUptime(w http.ResponseWriter, r *http.Request) {

	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	type dummyWorker struct {
		Url          string `json:"url"`
		Timezone     string `json:"timezone"`
		PollInterval int    `json:"pollInterval"`
	}
	var dummy dummyWorker
	if err := json.NewDecoder(r.Body).Decode(&dummy); err != nil || strings.TrimSpace(dummy.Url) == "" || dummy.PollInterval < 10 || dummy.PollInterval > 60 {
		http.Error(w, "Invalid worker URL", http.StatusBadRequest)
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	response, err := client.GetWorkerUptime(context.Background(), &pb.WorkerUrl{Url: dummy.Url, Timezone: dummy.Timezone, PollInterval: int32(dummy.PollInterval)})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	httpresponse, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(httpresponse)
}
