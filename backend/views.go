package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func welcome(w http.ResponseWriter, r *http.Request) {
	var jsonResponse = make(map[string]interface{})
	var unauthorized = false
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse["err"] = errors.New("Unauthorized")
		unauthorized = true
	}
	if !unauthorized {
		token := r.Header.Get("Authorization")
		jsonResponse["exists"] = token != ""
		w.WriteHeader(http.StatusOK)
	}
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
		if error_db := createUser(payload.Mail, payload.Secret, payload.AppName); error_db != nil {
			jsonResponse["status"] = false
			jsonResponse["err"] = "Error occured during account creation. Please try again later."
		} else {
			jsonResponse["status"] = true
		}
	} else {
		created := signIn(payload.Mail, payload.Secret)
		if created {
			jsonResponse["status"] = true
		} else {
			jsonResponse["status"] = false
			jsonResponse["err"] = "Invalid credentials"
		}
	}
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}
