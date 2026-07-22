package db

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var usersPool *pgxpool.Pool

func SetUsersPool(pool *pgxpool.Pool) {
	usersPool = pool
}
func getUsersPool() (*pgxpool.Pool, error) {
	if usersPool == nil {
		return nil, errors.New("users database pool is not initialized")
	}
	return usersPool, nil
}

var logsPool *pgxpool.Pool

func SetLogsPool(pool *pgxpool.Pool) {
	logsPool = pool
}

func GetLogsPool() (*pgxpool.Pool, error) {
	if logsPool == nil {
		return nil, errors.New("logs database pool is not initialized")
	}
	return logsPool, nil
}
