package api

import (
	"backend/auth"
	"backend/db"
	"context"
	"encoding/json"
	"net/http"
)

func AddCustomMetric(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	var metric db.CustomMetric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := metric.Initialize(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ListCustomMetrics(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	metrics, err := db.ListCustomMetrics(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var formattedMetrics []byte
	if formattedMetrics, err = json.Marshal(metrics); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(formattedMetrics)
}

func UpdateCustomMetrics(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	type payload struct {
		Original db.CustomMetric `json:"original"`
		Updated  db.CustomMetric `json:"updated"`
	}
	var data payload
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := data.Original.Edit(ctx, data.Updated); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteCustomMetric(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	var metric db.CustomMetric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := metric.Delete(ctx); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
