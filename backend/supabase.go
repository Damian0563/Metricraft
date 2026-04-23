package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func createUser(mail string, secret string, appName string) (string, error) {
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

func signIn(mail string, secret string) (string, bool) {
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

func getAppNameByToken(token string) (string, error) {
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

func getUserByToken(token string) (User, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return User{}, err
	}
	defer conn.Close(context.Background())
	var user User
	err = conn.QueryRow(context.Background(), "SELECT mail, app_name, uuid, realtime FROM users WHERE uuid = $1", token).Scan(&user.Mail, &user.AppName, &user.UUID, &user.Settings.Realtime)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func changeRealtimeByToken(token string, enabled bool) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	_, err = conn.Exec(context.Background(), "UPDATE users SET realtime = $1 WHERE uuid = $2", enabled, token)
	if err != nil {
		return err
	}
	return nil
}
