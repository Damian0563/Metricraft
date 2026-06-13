package api

import (
	"backend/db"
	"backend/mail"
	"backend/types"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"os"
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
	userId, err := db.FetchRecoveryUser(id)
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
		db.InvalidateSession(userId)
		db.InvalidateRecovery(id)
		w.WriteHeader(http.StatusOK)
	}
}

func SendRecovery(w http.ResponseWriter, r *http.Request) {
	if os.Getenv("SECRET") != r.Header.Get("Authorization") && os.Getenv("SECRET") != "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var payload types.SendRecoveryUser
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid payload", http.StatusInternalServerError)
		return
	}
	if user, err := db.GetRecoveryUser(payload.Mail); err != nil {
		if err.Error() == "User to be recovered not found." {
			http.Error(w, "User to be recovered not found.", http.StatusBadRequest)
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	} else {
		id := uuid.New().String()
		err = db.SetRecovery(user.Id, id)
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
	json.NewDecoder(r.Body).Decode(&payload)
	var routine = make(chan types.ExistsErrResponse, 1)
	mailAddress := payload.Mail
	go db.CheckUserExists(routine, mailAddress)
	response := <-routine
	if response.Err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if response.Exists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	code := db.GenerateCode()
	err := db.SetCodeValidity(mailAddress, code)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = mail.SendVerification(mailAddress, code)
	if err != nil {
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
	json.NewDecoder(r.Body).Decode(&payload)
	var routine = make(chan types.ExistsErrResponse, 2)
	go db.CheckAllowed(routine, payload.Mail, payload.AppName)
	go db.CheckCodeValidity(routine, payload.Mail, payload.Code)
	var (
		codeValid   bool
		permitted   bool
		internalErr bool
		malicious   bool
	)
	for received := 2; received > 0; received-- {
		select {
		case response := <-routine:
			switch response.Origin {
			case "checkCodeValidity":
				if response.Err != nil {
					internalErr = true
				} else {
					codeValid = response.Exists
				}
			case "checkAllowed":
				msg := ""
				if response.Err != nil {
					msg = response.Err.Error()
				}
				switch {
				case msg == "Owner is allowed to sign in", response.Err == nil && response.Exists:
					permitted = true
				case msg == "Permission needed from the owner", response.Err == nil && !response.Exists:
					if err := db.AddToPendingList(payload.Mail, payload.AppName); err != nil {
						internalErr = true
					} else {
						err := mail.SendPermissionRequest(response.Owner, payload.Mail, payload.AppName)
						if err != nil {
							internalErr = true
						} else {
							permitted = false
						}
					}
				case msg == "App name verification needed.", response.Err == nil && response.Exists:
					malicious = !db.VerifyAppName(payload.AppName)
				default:
					internalErr = true
				}
			}
		}
	}
	switch {
	case internalErr:
		http.Error(w, "Something went wrong, please try again later", http.StatusInternalServerError)
	case !codeValid:
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
	case !permitted:
		http.Error(w, "Permission needed from the owner", http.StatusUnauthorized)
	case malicious:
		http.Error(w, "Invalid app name", http.StatusForbidden)
	default:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Verification successful"))
	}
}
