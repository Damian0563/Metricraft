package api

import (
	"backend/auth"
	"backend/db"
	"backend/types"
	"context"
	"encoding/json"
	"net/http"
)

func AddRule(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	if authed := token.ValidateRequest(&w, true); !authed {
		return
	}
	var rule types.Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := db.AddRule(ctx, rule); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetRules(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	if authed := token.ValidateRequest(&w, false); !authed {
		return
	}
	ctx := context.Background()
	rules, err := db.GetRules(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(rules)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

func DeleteRule(w http.ResponseWriter, r *http.Request) {
	token := auth.NewToken(r.Header.Get("Session-Token"))
	if authed := token.ValidateRequest(&w, false); !authed {
		return
	}
	var rule types.Rule
	if err := json.NewDecoder(r.Body).Decode(&rule); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	if err := db.DeleteRule(ctx, rule); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
