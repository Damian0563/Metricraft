package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var jsonResponse = make(map[string]interface{})
	if _, err := os.Stat("exists.txt"); err == nil {
		jsonResponse["exists"] = true
	} else if errors.Is(err, os.ErrNotExist) {
		jsonResponse["exists"] = false
	} else {
		jsonResponse["err"] = err
	}
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}
