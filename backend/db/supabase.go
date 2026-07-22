package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"slices"
	"strings"
	"time"
)

func CheckUserExists(routine chan types.ExistsErrResponse, mail string) {
	ctx := context.Background()
	conn, err := getUsersPool()
	var origin string = "checkUserExists"
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	var exists bool
	err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE mail = $1)", mail).Scan(&exists)
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	routine <- types.ExistsErrResponse{Exists: exists, Err: err, Origin: origin, Owner: ""}
}

func CheckAllowed(routine chan types.ExistsErrResponse, mail string, appName string) {
	ctx := context.Background()
	conn, err := getUsersPool()
	var origin string = "checkAllowed"
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	var allowedUsers string
	var owner string
	err = conn.QueryRow(ctx, "SELECT allowed_users,mail FROM users WHERE app_name=$1", appName).Scan(&allowedUsers, &owner)
	if err != nil {
		if err == pgx.ErrNoRows {
			routine <- types.ExistsErrResponse{Exists: true, Err: errors.New("App name verification needed."), Origin: origin, Owner: owner}
			return
		}
	}
	if owner == mail {
		routine <- types.ExistsErrResponse{Exists: true, Err: errors.New("Owner is allowed to sign in"), Origin: origin, Owner: owner}
	} else if allowedUsers == "" && owner == "" {
		routine <- types.ExistsErrResponse{Exists: false, Err: nil, Origin: origin, Owner: owner}
	} else if allowedUsers != "" && owner != "" {
		var allowed []types.AllowedUsersDb
		err = json.Unmarshal([]byte(allowedUsers), &allowed)
		if err != nil {
			routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: owner}
			return
		}
		for _, user := range allowed {
			if user.Mail == mail {
				routine <- types.ExistsErrResponse{Exists: true, Err: nil, Origin: origin, Owner: owner}
				return
			}
		}
		routine <- types.ExistsErrResponse{Exists: false, Err: errors.New("Permission needed from the owner"), Origin: origin, Owner: owner}
	} else {
		routine <- types.ExistsErrResponse{Exists: false, Err: errors.New("Something went wrong"), Origin: origin, Owner: owner}
	}
}

func GetRecoveryUser(recoveryMail string) (types.SendRecoveryUser, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return types.SendRecoveryUser{}, err
	}
	var id string
	var mail string
	err = conn.QueryRow(ctx, "SELECT uuid,mail FROM users WHERE mail=$1", recoveryMail).Scan(&id, &mail)
	if err != nil {
		return types.SendRecoveryUser{}, err
	} else if id == "" {
		return types.SendRecoveryUser{}, errors.New("User to be recovered not found.")
	}
	return types.SendRecoveryUser{Mail: mail, Id: id}, nil
}

func ChangePassword(userId string, password string) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "UPDATE users SET secret=$1 WHERE uuid=$2", string(hashedSecret), userId)
	return err
}

func AllowUsers(mails []string, appName string, errChannel chan error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		errChannel <- err
		return
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		errChannel <- err
		return
	}
	defer tx.Rollback(ctx)
	var allowedUsers string
	err = tx.QueryRow(ctx, "SELECT allowed_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&allowedUsers)
	if err != nil {
		errChannel <- err
		return
	}
	var allowed []types.AllowedUsersDb
	if err = json.Unmarshal([]byte(allowedUsers), &allowed); err != nil {
		errChannel <- err
		return
	}
	for _, mail := range mails {
		found := false
		for _, user := range allowed {
			if user.Mail == mail {
				found = true
				break
			}
		}
		if !found {
			allowed = append(allowed, types.AllowedUsersDb{Mail: mail, ReceiveNotifications: false})
		}
	}
	var marshaled []byte
	if marshaled, err = json.Marshal(allowed); err != nil {
		errChannel <- err
		return
	}
	_, err = tx.Exec(ctx, "UPDATE users SET allowed_users=$1 WHERE app_name=$2 AND owner=true", string(marshaled), appName)
	if err != nil {
		errChannel <- err
		return
	}
	errChannel <- tx.Commit(ctx)
}

func AddToPendingList(mail string, appName string) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	var pendingUsers string
	err = conn.QueryRow(ctx, "SELECT pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&pendingUsers)
	if err != nil {
		return err
	}
	var pending []string
	err = json.Unmarshal([]byte(pendingUsers), &pending)
	if err != nil {
		return err
	}
	pending = append(pending, mail)
	if err != nil {
		return err
	}
	_, err = conn.Exec(ctx, "UPDATE users SET pending_users=$1 WHERE app_name=$2 AND owner=true", pending, appName)
	return err
}

func unMarshalUsers(jsonString string, jsonSlice *[]string) {
	err := json.Unmarshal([]byte(jsonString), jsonSlice)
	if err != nil {
		return
	}
}

func ChangeNotificationRecipients(allowedUsers []types.AllowedUsersDb, appName string) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	if _, err = tx.Exec(ctx, "SELECT allowed_users FROM users WHERE app_name=$1 AND owner=true FOR UPDATE", appName); err != nil && err != pgx.ErrNoRows {
		return err
	}
	var marshaled []byte
	if marshaled, err = json.Marshal(allowedUsers); err != nil {
		return err
	}
	if _, err = tx.Exec(ctx, "UPDATE users SET allowed_users=$1 WHERE app_name=$2 AND owner=true", string(marshaled), appName); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func HandleInvite(mail string, decision string, appName string) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	var pendingUsers string
	var allowedUsers string
	err = tx.QueryRow(ctx, "SELECT allowed_users,pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&allowedUsers, &pendingUsers)
	if err != nil {
		return err
	}
	var pending []string
	var allowed []types.AllowedUsersDb
	unMarshalUsers(pendingUsers, &pending)
	if err = json.Unmarshal([]byte(allowedUsers), &allowed); err != nil {
		return err
	}
	idx := slices.Index(pending, mail)
	if idx == -1 {
		return errors.New("User not found")
	}
	pending = slices.Delete(pending, idx, idx+1)
	if decision == "true" {
		allowed = append(allowed, types.AllowedUsersDb{Mail: mail, ReceiveNotifications: false})
	}
	var marshaled []byte
	if marshaled, err = json.Marshal(allowed); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "UPDATE users SET pending_users=$1,allowed_users=$2 WHERE app_name=$3 AND owner=true", pending, marshaled, appName)
	return tx.Commit(ctx)
}

func GetTeamUsers(appName string) ([]types.AllowedUsers, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return nil, err
	}
	var teamUsers string
	var owner string
	err = conn.QueryRow(ctx, "SELECT mail,allowed_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&owner, &teamUsers)
	if err != nil {
		return nil, err
	}
	var users []types.AllowedUsersDb
	err = json.Unmarshal([]byte(teamUsers), &users)
	if err != nil {
		return nil, err
	}
	var response []types.AllowedUsers
	for _, user := range users {
		var exists bool
		err = conn.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE mail = $1)", user.Mail).Scan(&exists)
		if err != nil {
			return nil, err
		}
		response = append(response, types.AllowedUsers{Mail: user.Mail, Initials: strings.ToUpper(user.Mail[0:2]), Status: exists, ReceiveNotifications: user.ReceiveNotifications})
	}
	return response, nil
}

func GetPendingUsers(appName string) ([]string, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return nil, err
	}
	var pendingUsers string
	err = conn.QueryRow(ctx, "SELECT pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&pendingUsers)
	if err != nil {
		return nil, err
	}
	var users []string
	err = json.Unmarshal([]byte(pendingUsers), &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func CreateUser(mail string, secret string, appName string) (string, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return "", err
	}
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	uuid := uuid.New().String()
	owner := checkAppNameExists(appName)
	allowedUsers, err := json.Marshal([]types.AllowedUsersDb{{
		Mail:                 mail,
		ReceiveNotifications: true,
	}})
	if err != nil {
		return "", err
	}
	if _, err = conn.Exec(ctx, "INSERT INTO users (created_at, mail, secret, app_name, uuid,owner,allowed_users) VALUES ($1, $2, $3, $4, $5,$6, $7)", time.Now(), mail, string(hashedSecret), appName, uuid, owner, string(allowedUsers)); err != nil {
		return "", err
	}
	if _, err = conn.Exec(ctx, "INSERT INTO workers (app_name) VALUES ($1)", appName); err != nil {
		return "", err
	}
	return uuid, nil
}

func checkAppNameExists(appName string) bool {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return false
	}
	var exists bool
	err = conn.QueryRow(ctx, "SELECT NOT EXISTS(SELECT 1 FROM users WHERE app_name = $1)", appName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func SignIn(mail string, secret string) (string, bool) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return "", false
	}
	var hashedSecret string
	var uuid string
	err = conn.QueryRow(ctx, "SELECT uuid, secret FROM users WHERE mail = $1", mail).Scan(&uuid, &hashedSecret)
	if err != nil {
		return "", false
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret))
	return uuid, err == nil
}

func GetAppNameByToken(token string) (string, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return "", err
	}
	var appName string
	err = conn.QueryRow(ctx, "SELECT app_name FROM users WHERE uuid = $1", token).Scan(&appName)
	if err != nil {
		return "", err
	}
	return appName, nil
}
