package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

// var DATABASE_URL string = "postgres://postgres:password@postgresql:5432/postgres?sslmode=disable"
var DATABASE_URL string = "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"

func initDB(ctx context.Context, errChannel chan error) {
	time.Sleep(15 * time.Second)
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		errChannel <- err
		return
	}
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS settings (id UUID PRIMARY KEY, realtime BOOL, enabled TEXT NOT NULL)")
	if err != nil {
		errChannel <- err
		return
	}
	_, err = conn.Exec(context.Background(), "CREATE TABLE IF NOT EXISTS logs (id UUID PRIMARY KEY,date TIMESTAMP,responseTime TEXT NOT NULL)")
	if err != nil {
		errChannel <- err
		return
	}
	errChannel <- conn.Close(context.Background())
}
