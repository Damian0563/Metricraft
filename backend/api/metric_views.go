package api

import (
	"backend/auth"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var grpcConn *grpc.ClientConn

func SetGRPCConn(conn *grpc.ClientConn) {
	grpcConn = conn
}

type resolutionDays struct {
	Days int32
}

func convertTimeframe(timeframe string) (time.Time, resolutionDays) {
	if timeframe == "" {
		timeframe = "7d"
	}
	timeframe = strings.ReplaceAll(timeframe, "d", "")
	num, err := strconv.Atoi(timeframe)
	if err != nil {
		num = 7
		fmt.Println(err)
	}
	now := time.Now()
	var timeResolution map[int]resolutionDays = map[int]resolutionDays{
		7:   {Days: 1},
		14:  {Days: 2},
		30:  {Days: 3},
		90:  {Days: 7},
		180: {Days: 14},
		365: {Days: 30},
	}
	return now.Add(-time.Hour * 24 * time.Duration(num)), timeResolution[num]
}

func Navigator(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed, err := token.Verify()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if grpcConn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(grpcConn)
	metric := r.URL.Query().Get("metric")
	timeframe := r.URL.Query().Get("timeframe")
	persist := r.URL.Query().Get("persist")
	persistChan := make(chan error, 1)
	if persist == "true" {
		go persistTimeframeSelection(persistChan, metric, timeframe)
	} else {
		persistChan <- nil
	}
	convertedTimeframe, resolution := convertTimeframe(timeframe)
	var response any
	switch metric {
	case "Traffic congestion trends":
		response, err = client.GetTrafficCongestion(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days})
	case "Geographical traffic":
		response, err = client.GetGeographicalTraffic(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days})
	case "P95 Latency":
		response, err = client.GetP95Latency(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days})
	case "Uptime Score":
		response, err = client.GetUptimeScore(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days})
	case "Throughput":
		response, err = client.GetThroughput(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe), Resolution: resolution.Days})
	default:
		break
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = <-persistChan //this can be swallowed internally, even if error occured
	w.Header().Set("Content-Type", "application/json")
	httpresponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(httpresponse)
}
