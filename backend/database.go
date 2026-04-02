package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"os"
	"time"
)

func createUser(mail string, secret string, appName string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO users (created_at, mail, secret, app_name) VALUES ($1, $2, $3, $4)", time.Now(), mail, string(hashedSecret), appName)
	if err != nil {
		return err
	}
	return nil
}

func signIn(mail string, secret string) bool {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	var hashedSecret string
	err = conn.QueryRow(context.Background(), "SELECT secret FROM users WHERE mail = $1", mail).Scan(&hashedSecret)
	if err != nil {
		panic(err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedSecret), []byte(secret))
	return err == nil
}
