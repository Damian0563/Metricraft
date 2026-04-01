package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

func sign(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	var payload signPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	var jsonResponse = make(map[string]interface{})
	if payload.AppName != "" {
		createUser(payload.Mail, payload.Secret, payload.AppName)
		jsonResponse["status"] = true
	} else {
		signIn(payload.Mail, payload.Secret)
		jsonResponse["status"] = true
	}
	fmt.Println(jsonResponse)
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}
