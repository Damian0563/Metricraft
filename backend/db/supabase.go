package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"slices"
	"strings"
	"time"
)

func CheckUserExists(routine chan types.ExistsErrResponse, mail string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	var origin string = "checkUserExists"
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	defer conn.Close(context.Background())
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE mail = $1)", mail).Scan(&exists)
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	routine <- types.ExistsErrResponse{Exists: exists, Err: err, Origin: origin, Owner: ""}
}

func CheckAllowed(routine chan types.ExistsErrResponse, mail string, appName string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	var origin string = "checkAllowed"
	if err != nil {
		routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: ""}
		return
	}
	defer conn.Close(context.Background())
	var allowed_users string
	var owner string
	err = conn.QueryRow(context.Background(), "SELECT allowed_users,mail FROM users WHERE app_name=$1", appName).Scan(&allowed_users, &owner)
	if err != nil {
		if err == pgx.ErrNoRows {
			routine <- types.ExistsErrResponse{Exists: true, Err: errors.New("App name verification needed."), Origin: origin, Owner: owner}
			return
		}
	}
	if owner == mail {
		routine <- types.ExistsErrResponse{Exists: true, Err: errors.New("Owner is allowed to sign in"), Origin: origin, Owner: owner}
	} else if allowed_users == "" && owner == "" {
		routine <- types.ExistsErrResponse{Exists: false, Err: nil, Origin: origin, Owner: owner}
	} else if allowed_users != "" && owner != "" {
		var allowed []string
		err = json.Unmarshal([]byte(allowed_users), &allowed)
		if err != nil {
			routine <- types.ExistsErrResponse{Exists: false, Err: err, Origin: origin, Owner: owner}
			return
		}
		for _, user := range allowed {
			if user == mail {
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
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return types.SendRecoveryUser{}, err
	}
	defer conn.Close(context.Background())
	var id string
	var mail string
	err = conn.QueryRow(context.Background(), "SELECT uuid,mail FROM users WHERE mail=$1", recoveryMail).Scan(&id, &mail)
	if err != nil {
		return types.SendRecoveryUser{}, err
	} else if id == "" {
		return types.SendRecoveryUser{}, errors.New("User to be recovered not found.")
	}
	return types.SendRecoveryUser{Mail: mail, Id: id}, nil
}

func ChangePassword(userId string, password string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), "UPDATE users SET secret=$1 WHERE uuid=$2", string(hashedSecret), userId)
	return err
}

func AddToPendingList(mail string, appName string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	var pendingUsers string
	err = conn.QueryRow(context.Background(), "SELECT pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&pendingUsers)
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
	_, err = conn.Exec(context.Background(), "UPDATE users SET pending_users=$1 WHERE app_name=$2 AND owner=true", pending, appName)
	return err
}

func unMarshalUsers(jsonString string, jsonSlice *[]string) {
	err := json.Unmarshal([]byte(jsonString), jsonSlice)
	if err != nil {
		return
	}
}

func HandleInvite(mail string, decision string, appName string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	var pendingUsers string
	var allowedUsers string
	err = conn.QueryRow(context.Background(), "SELECT allowed_users,pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&allowedUsers, &pendingUsers)
	if err != nil {
		return err
	}
	var pending []string
	var allowed []string
	unMarshalUsers(pendingUsers, &pending)
	unMarshalUsers(allowedUsers, &allowed)
	idx := slices.Index(pending, mail)
	if idx == -1 {
		return errors.New("User not found")
	}
	pending = slices.Delete(pending, idx, idx+1)
	if decision == "true" {
		allowed = append(allowed, mail)
	}
	_, err = conn.Exec(context.Background(), "UPDATE users SET pending_users=$1,allowed_users=$2 WHERE app_name=$3 AND owner=true", pending, allowed, appName)
	return err
}

func GetTeamUsers(appName string) ([]types.AllowedUsers, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	var teamUsers string
	var owner string
	err = conn.QueryRow(context.Background(), "SELECT mail,allowed_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&owner, &teamUsers)
	if err != nil {
		return nil, err
	}
	var users []string
	err = json.Unmarshal([]byte(teamUsers), &users)
	if err != nil {
		return nil, err
	}
	users = slices.Insert(users, 0, owner)
	var response []types.AllowedUsers
	for _, user := range users {
		response = append(response, types.AllowedUsers{Mail: user, Initials: strings.ToUpper(user[0:2])})
	}
	return response, nil
}

func GetPendingUsers(appName string) ([]string, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	var pendingUsers string
	err = conn.QueryRow(context.Background(), "SELECT pending_users FROM users WHERE app_name=$1 AND owner=true", appName).Scan(&pendingUsers)
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
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	uuid := uuid.New().String()
	owner := checkAppNameExists(appName)
	_, err = conn.Exec(context.Background(), "INSERT INTO users (created_at, mail, secret, app_name, uuid,owner) VALUES ($1, $2, $3, $4, $5,$6)", time.Now(), mail, string(hashedSecret), appName, uuid, owner)
	if err != nil {
		return "", err
	}
	return uuid, nil
}

func checkAppNameExists(appName string) bool {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return false
	}
	defer conn.Close(context.Background())
	var exists bool
	err = conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM settings WHERE appname = $1)", appName).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func SignIn(mail string, secret string) (string, bool) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	var hashedSecret string
	var uuid string
	err = conn.QueryRow(context.Background(), "SELECT uuid, secret FROM users WHERE mail = $1", mail).Scan(&uuid, &hashedSecret)
	if err != nil {
		return "", false
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret))
	return uuid, err == nil
}

func GetAppNameByToken(token string) (string, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return "", err
	}
	defer conn.Close(context.Background())
	var appName string
	err = conn.QueryRow(context.Background(), "SELECT app_name FROM users WHERE uuid = $1", token).Scan(&appName)
	if err != nil {
		return "", err
	}
	return appName, nil
}
