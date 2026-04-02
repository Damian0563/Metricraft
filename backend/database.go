package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func createUser(mail string, secret string, appName string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO users (mail, secret, app_name) VALUES ($1, $2, $3)", mail, string(hashedSecret), appName)
	if err != nil {
		panic(err)
	}
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
