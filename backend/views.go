package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
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
		token := r.Header.Get("Session-Token")
		exists, err := checkSecret(token)
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

func dashboardInit(w http.ResponseWriter, r *http.Request) {
	signedToken := r.Header.Get("Session-Token")
	token := strings.Split(signedToken, ":")[0]
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
		signed, error := signSecret(token)
		if error != nil {
			Response.Error = "Error occured during signing. Please try again later."
			Response.SignedSecret = ""
		} else {
			Response.SignedSecret = signed
		}
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
	var errChannel = make(chan error)
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
			context := context.Background()
			go initDB(context, errChannel)
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
	err = <-errChannel
	var response []byte
	if err == nil {
		w.WriteHeader(http.StatusOK)
		response, _ = json.Marshal(jsonResponse)
	} else {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		jsonResponse["token"] = ""
		jsonResponse["err"] = "Error occured during account creation. Please try again later."
		response, _ = json.Marshal(jsonResponse)
	}
	w.Write(response)
}
