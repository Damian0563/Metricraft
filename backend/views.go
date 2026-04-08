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
		token := r.Header.Get("Session-Token")
		jsonResponse["exists"] = token != ""
		w.WriteHeader(http.StatusOK)
	}
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}

func dashboardInit(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Session-Token")
	var Response = dashboardInitPayload{}
	appName, err := getAppNameByToken(token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Response.Error = err.Error()
		Response.AppName = ""
		Response.SignedSecret = ""
	} else {
		w.WriteHeader(http.StatusOK)
		Response.AppName = appName
		Response.SignedSecret = token //for now
		Response.Error = ""
	}
	response, _ := json.Marshal(Response)
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
		if uuid, error_db := createUser(payload.Mail, payload.Secret, payload.AppName); error_db != nil {
			jsonResponse["token"] = ""
			jsonResponse["err"] = "Error occured during account creation. Please try again later."
		} else {
			jsonResponse["token"] = uuid
		}
	} else {
		uuid, ok := signIn(payload.Mail, payload.Secret)
		if ok {
			jsonResponse["token"] = uuid
		} else {
			jsonResponse["token"] = ""
			jsonResponse["err"] = "Invalid credentials"
		}
	}
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}
