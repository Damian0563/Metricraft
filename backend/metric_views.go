package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func getCongestion(w http.ResponseWriter, r *http.Request) {
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
	type congestionPayload struct {
		Timeframe string `json:"timeframe"`
	}
	var payload congestionPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	timeframe := payload.Timeframe
	fmt.Println(timeframe)
	w.WriteHeader(http.StatusOK)
}
