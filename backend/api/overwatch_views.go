package api

import (
	"backend/auth"
	"backend/db"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
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
	tz := strings.TrimSpace(r.URL.Query().Get("timezone"))
	ctx := context.Background()
	if err := metric.Initialize(ctx, tz); err != nil {
		if errors.Is(err, db.ErrInvalidTimezone) {
			http.Error(w, "Invalid timezone", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
