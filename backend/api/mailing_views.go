package api

import (
	"backend/db"
	"backend/mail"
	"backend/redis"
	"backend/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
	"strings"
)

func CheckRecovery(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userId, err := redis.FetchRecoveryUser(id)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	var payload types.RecoveryCheckPayload
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = db.ChangePassword(userId, payload.Password); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		redis.InvalidateSession(userId)
		redis.InvalidateRecovery(id)
		w.WriteHeader(http.StatusOK)
	}
}

func SendRecovery(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var payload types.SendRecoveryUser
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !mail.ValidateMail(payload.Mail) {
		http.Error(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	payload.Mail = strings.TrimSpace(payload.Mail)
	if user, err := db.GetRecoveryUser(payload.Mail); err != nil {
		if err.Error() == "User to be recovered not found." {
			http.Error(w, "User to be recovered not found.", http.StatusBadRequest)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	} else {
		id := uuid.New().String()
		err = redis.SetRecovery(user.Id, id)
		if err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		if base := os.Getenv("ALLOWED_ORIGINS"); base != "" {
			link := fmt.Sprintf("%s/recovery/%s", base, id)
			err = mail.SendRecovery(user.Mail, "Password Recovery", link)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func SendVerification(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type sendVerificationPayload struct {
		Mail string `json:"mail"`
	}
	var payload sendVerificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !mail.ValidateMail(payload.Mail) {
		http.Error(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	var routine = make(chan types.ExistsErrResponse, 1)
	mailAddress := strings.TrimSpace(payload.Mail)
	go db.CheckUserExists(routine, mailAddress)
	response := <-routine
	if response.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if response.Exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := redis.GenerateCode()
	if err := redis.SetCodeValidity(mailAddress, code); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := mail.SendVerification(mailAddress, code); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckVerification(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	type checkVerificationPayload struct {
		AppName string `json:"appName"`
		Mail    string `json:"mail"`
		Code    string `json:"code"`
	}
	var payload checkVerificationPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}
	if !mail.ValidateMail(payload.Mail) {
		http.Error(w, "Invalid email address.", http.StatusBadRequest)
		return
	}
	payload.Mail = strings.TrimSpace(payload.Mail)
	routine := make(chan types.ExistsErrResponse, 1)
	go redis.CheckCodeValidity(routine, payload.Mail, payload.Code)
	codeResponse := <-routine
	if codeResponse.Err != nil {
		http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
		return
	}
	if !codeResponse.Exists {
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
		return
	}

	go db.CheckAllowed(routine, payload.Mail, payload.AppName)
	allowedResponse := <-routine
	msg := ""
	if allowedResponse.Err != nil {
		msg = allowedResponse.Err.Error()
	}
	switch {
	case msg == "Owner is allowed to sign in", allowedResponse.Err == nil && allowedResponse.Exists:
		w.WriteHeader(http.StatusOK)
	case msg == "Permission needed from the owner", allowedResponse.Err == nil && !allowedResponse.Exists:
		if err := db.AddToPendingList(payload.Mail, payload.AppName); err != nil {
			http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
			return
		}
		if err := mail.SendPermissionRequest(allowedResponse.Owner, payload.Mail, payload.AppName); err != nil {
			http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
			return
		}
		http.Error(w, "Permission needed from the owner", http.StatusUnauthorized)
	case msg == "App name verification needed.", allowedResponse.Err == nil && allowedResponse.Exists:
		if !db.VerifyAppName(payload.AppName) {
			http.Error(w, "Invalid app name", http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
	default:
		http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
	}
}

func SendErrorNotification(ctx context.Context, appName string, url string, statusCode int, errMsg string, errChan chan error) {
	teamUsers, err := db.GetTeamUsers(appName)
	if err != nil {
		errChan <- err
		return
	}
	var recipients []string
	for _, user := range teamUsers {
		if user.ReceiveNotifications {
			recipients = append(recipients, user.Mail)
		}
	}
	errChan <- mail.SendErrorNotification(recipients, appName, url, statusCode, errMsg)
}
