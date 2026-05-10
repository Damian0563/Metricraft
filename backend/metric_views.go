package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

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
	metric := r.URL.Query().Get("metric")
	timeframe := r.URL.Query().Get("timeframe")
	var response any
	switch metric {
	case "Traffic congestion trends":
		response, err = getCongestion(timeframe)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
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

func getCongestion(timeframe string) (map[string]int32, error) {
	return nil, nil
}
