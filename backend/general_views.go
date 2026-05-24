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
		jsonResponse["exists"] = false
		unauthorized = true
	}
	if !unauthorized {
		token := Token{token: r.Header.Get("Session-Token")}
		exists, err := token.verify()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse["err"] = "Error occured during checking session token. Please try again later."
			jsonResponse["exists"] = false
			return
		}
		jsonResponse["exists"] = exists
		w.WriteHeader(http.StatusOK)
	}
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}

func toggleRealtime(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := Token{token: r.Header.Get("Session-Token")}
	authed := token.validateRequest(&w, true)
	if !authed {
		return
	}
	type realtimePayload struct {
		Enabled bool `json:"enabled"`
	}
	var payload realtimePayload
	json.NewDecoder(r.Body).Decode(&payload)
	err := ChangeRealtime(payload.Enabled)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func changeMetrics(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := Token{token: r.Header.Get("Session-Token")}
	authed := token.validateRequest(&w, true)
	if !authed {
		return
	}
	var payload []Metric
	json.NewDecoder(r.Body).Decode(&payload)
	err := ChangeMetrics(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func changeRetention(w http.ResponseWriter, r *http.Request) {
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
	} else {
		if err = token.updateToken(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	type retentionPayload struct {
		Retention int `json:"retention"`
	}
	var payload retentionPayload
	json.NewDecoder(r.Body).Decode(&payload)
	err = ChangeLogsRetention(payload.Retention)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func dashboardInit(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := Token{token: r.Header.Get("Session-Token")}
	authed := token.validateRequest(&w, false)
	if !authed {
		return
	}
	var Response = dashboardInitPayload{}
	Response.SignedSecret = ""
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Response.Error = err.Error()
		Response.AppName = ""
	} else {
		Response.AppName = appName
		signed, error := token.sign(false)
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Response.Error = "Error occured during signing. Please try again later."
		} else {
			status := http.StatusOK
			Response.SignedSecret = signed
			Response.Settings, err = GetSettings()
			if err != nil {
				Response.Error = "Error occured during fetching settings. Please try again later."
				status = http.StatusInternalServerError
			}
			Response.Urls, err = GetUrls()
			if err != nil {
				Response.Error = "Error occured during fetching urls. Please try again later."
				status = http.StatusInternalServerError
			}
			w.WriteHeader(status)
		}
		response, _ := json.Marshal(Response)
		w.Write(response)
	}
}

func sign(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	type signPayload struct {
		Mail    string `json:"mail"`
		Secret  string `json:"secret"`
		AppName string `json:"appName",omitempty`
	}
	var payload signPayload
	err = json.Unmarshal(body, &payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	var jsonResponse = make(map[string]interface{})
	if payload.AppName != "" {
		if uuid, error_db := createUser(payload.Mail, payload.Secret, payload.AppName); error_db != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse["token"] = ""
			jsonResponse["err"] = "Error occured during account creation. Please try again later."
		} else {
			w.WriteHeader(http.StatusOK)
			token := Token{token: uuid}
			signed, err := token.sign(true)
			if err != nil {
				jsonResponse["token"] = ""
				jsonResponse["err"] = "Error occured during signing. Please try again later."
			} else {
				jsonResponse["token"] = signed
			}
		}
	} else {
		uuid, ok := signIn(payload.Mail, payload.Secret)
		if ok {
			token := Token{token: uuid}
			signed, err := token.sign(true)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				jsonResponse["token"] = ""
				jsonResponse["err"] = "Error occured during signing. Please try again later."
			}
			w.WriteHeader(http.StatusOK)
			jsonResponse["token"] = signed
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			jsonResponse["token"] = ""
			jsonResponse["err"] = "Invalid credentials"
		}
	}
	response, _ := json.Marshal(jsonResponse)
	w.Write(response)
}
