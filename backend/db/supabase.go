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
	_, err = conn.Exec(context.Background(), "INSERT INTO users (created_at, mail, secret, app_name, uuid) VALUES ($1, $2, $3, $4, $5)", time.Now(), mail, string(hashedSecret), appName, uuid)
	if err != nil {
		return "", err
	}
	return uuid, nil
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
