package main

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func convertTimeframe(timeframe string) time.Time {
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
	return now.Add(-time.Hour * 24 * time.Duration(num))
}

func navigator(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := Token{token: r.Header.Get("Session-Token")}
	authed, err := token.verify()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if conn == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client := pb.NewMetricraftClient(conn)
	metric := r.URL.Query().Get("metric")
	timeframe := r.URL.Query().Get("timeframe")
	convertedTimeframe := convertTimeframe(timeframe)
	fmt.Println(metric)
	var response any
	switch metric {
	case "Traffic congestion trends":
		response, err = client.GetTrafficCongestion(context.Background(), &pb.Timeframe{Start: timestamppb.New(convertedTimeframe)})
		fmt.Println(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Println(response)
	default:
		break
	}
	w.Header().Set("Content-Type", "application/json")
	httpresponse, err := json.Marshal(response)
	w.Write(httpresponse)
}
