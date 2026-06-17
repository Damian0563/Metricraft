package api

import (
	"backend/auth"
	"backend/db"
	mailer "backend/mail"
	"backend/types"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	var jsonResponse = make(map[string]interface{})
	var unauthorized = false
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		w.WriteHeader(http.StatusUnauthorized)
		jsonResponse["err"] = errors.New("Unauthorized")
		jsonResponse["exists"] = false
		unauthorized = true
	}
	if !unauthorized {
		token := auth.NewToken(r.Header.Get("Session-Token"))
		exists, err := token.Verify()
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

func ToggleRealtime(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
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

func ChangeMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	var payload []types.Metric
	json.NewDecoder(r.Body).Decode(&payload)
	err := ChangeMetrics(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func ChangeRetention(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed, err := token.Verify()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		if err = token.UpdateToken(); err != nil {
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

func TeamMembers(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type teamMembersPayload struct {
		Users []types.AllowedUsers `json:"users"`
	}
	var payload teamMembersPayload
	payload.Users, err = db.GetTeamUsers(appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func SendInvites(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else if mode == "" {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if mode == "manual" {
		type Invite struct {
			Invitees []string `json:"invitees"`
		}
		var invitees Invite
		json.NewDecoder(r.Body).Decode(&invitees)
		err := mailer.SendManualInvites(invitees.Invitees)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else if mode == "batch" {

	} else {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
	}
}

func HandleInvite(w http.ResponseWriter, r *http.Request) {
	mail := r.URL.Query().Get("user")
	decision := r.URL.Query().Get("action")
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	} else if mail == "" || decision == "" {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = db.HandleInvite(mail, decision, appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = mailer.NotifyDecision(mail, decision, appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PendingInvites(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		return
	}
	appName, err := token.GetAppName()
	pendingUsers, err := db.GetPendingUsers(appName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var jsonResponse = make(map[string][]types.PendingUsers)
	var formattedUsers []types.PendingUsers
	for _, user := range pendingUsers {
		formattedUsers = append(formattedUsers, types.PendingUsers{Mail: user})
	}
	jsonResponse["users"] = formattedUsers
	response, err := json.Marshal(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func DashboardInit(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		return
	}
	var Response = types.DashboardInitPayload{}
	Response.SignedSecret = ""
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Response.Error = err.Error()
		Response.AppName = ""
	} else {
		Response.AppName = appName
		signed, error := token.Sign(false)
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

func Sign(w http.ResponseWriter, r *http.Request) {
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
		if uuid, error_db := db.CreateUser(payload.Mail, payload.Secret, payload.AppName); error_db != nil {
			w.WriteHeader(http.StatusInternalServerError)
			jsonResponse["token"] = ""
			jsonResponse["err"] = "Error occured during account creation. Please try again later."
		} else {
			w.WriteHeader(http.StatusOK)
			token := auth.NewToken(uuid)
			signed, err := token.Sign(true)
			if err != nil {
				jsonResponse["token"] = ""
				jsonResponse["err"] = "Error occured during signing. Please try again later."
			} else {
				jsonResponse["token"] = signed
			}
		}
	} else {
		uuid, ok := db.SignIn(payload.Mail, payload.Secret)
		if ok {
			token := auth.NewToken(uuid)
			signed, err := token.Sign(true)
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
