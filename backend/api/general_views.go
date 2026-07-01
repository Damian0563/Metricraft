package api

import (
	"backend/auth"
	"backend/db"
	mailer "backend/mail"
	"backend/types"
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

func Welcome(w http.ResponseWriter, r *http.Request) {
	var jsonResponse = make(map[string]any)
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	type realtimePayload struct {
		Enabled bool `json:"enabled"`
	}
	var payload realtimePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := ChangeRealtime(payload.Enabled); err != nil {
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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	var payload []types.Metric
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := ChangeMetrics(payload); err != nil {
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
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = ChangeLogsRetention(payload.Retention); err != nil {
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
		w.WriteHeader(http.StatusUnauthorized)
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

func UploadUsersFromCSV(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if file, err := csv.NewReader(r.Body).ReadAll(); err == nil {
		if len(file) <= 1 {
			http.Error(w, "No mails found, expected structure is one column 'emails' with emails enlisted below in csv format.", http.StatusBadRequest)
			return
		} else if len(file) > 51 {
			http.Error(w, "Too many mails, maximum is 50.", http.StatusBadRequest)
			return
		}
		var invitees types.Invite
		for _, row := range file[1:] {
			address := strings.TrimSpace(row[0])
			if !mailer.ValidateMail(address) {
				http.Error(w, "Invalid email address: "+address, http.StatusBadRequest)
				return
			}
			invitees.Invitees = append(invitees.Invitees, address)
		}
		if err = sendInvites(invitees.Invitees, appName); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func sendInvites(invitees []string, appName string) error {
	errChannel := make(chan error, 2)
	go mailer.SendManualInvites(invitees, appName, errChannel)
	go db.AllowUsers(invitees, appName, errChannel)
	var inviteErr error
	for range 2 {
		if err := <-errChannel; err != nil && inviteErr == nil {
			inviteErr = err
		}
	}
	return inviteErr
}

func SendInvites(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, true)
	if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var invitees types.Invite
	if err = json.NewDecoder(r.Body).Decode(&invitees); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	for i, address := range invitees.Invitees {
		address = strings.TrimSpace(address)
		if !mailer.ValidateMail(address) {
			http.Error(w, "Invalid email address: "+address, http.StatusBadRequest)
			return
		}
		invitees.Invitees[i] = address
	}
	if err = sendInvites(invitees.Invitees, appName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
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
	} else if !mailer.ValidateMail(mail) {
		http.Error(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	mail = strings.TrimSpace(mail)
	token := auth.NewToken(r.Header.Get("Session-Token"))
	authed := token.ValidateRequest(&w, false)
	if !authed {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	appName, err := token.GetAppName()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = db.HandleInvite(mail, decision, appName); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err = mailer.NotifyDecision(mail, decision, appName); err != nil {
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
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		w.WriteHeader(http.StatusUnauthorized)
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
		signed, err := token.Sign(false)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			Response.Error = "Error occured during signing. Please try again later."
		} else {
			status := http.StatusOK
			Response.SignedSecret = signed
			if Response.Settings, err = GetSettings(); err != nil {
				Response.Error = "Error occured during fetching settings. Please try again later."
				status = http.StatusInternalServerError
			}
			if Response.Urls, err = GetUrls(); err != nil {
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
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	type signPayload struct {
		Mail    string `json:"mail"`
		Secret  string `json:"secret"`
		AppName string `json:"appName",omitempty`
	}
	var payload signPayload
	if err = json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !mailer.ValidateMail(payload.Mail) {
		http.Error(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	payload.Mail = strings.TrimSpace(payload.Mail)
	var jsonResponse = make(map[string]any)
	if payload.AppName != "" {
		if uuid, err := db.CreateUser(payload.Mail, payload.Secret, payload.AppName); err != nil {
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
		if uuid, ok := db.SignIn(payload.Mail, payload.Secret); ok {
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
