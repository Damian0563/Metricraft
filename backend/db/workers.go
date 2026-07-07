package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"os"
	"slices"
)

func SaveWorker(appName string, worker types.Worker, errChan chan error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		errChan <- err
		return
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		errChan <- err
		return
	}
	defer conn.Close(context.Background())
	var workers string
	err = conn.QueryRow(context.Background(), "SELECT workers FROM workers WHERE app_name=$1", appName).Scan(&workers)
	if err != nil {
		errChan <- err
		return
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		errChan <- err
		return
	}
	if len(workerList) >= 5 {
		errChan <- errors.New("Worker limit reached.")
		return
	}
	var marshalledWorker []byte
	if marshalledWorker, err = json.Marshal(worker); err != nil {
		errChan <- err
		return
	}
	for _, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			errChan <- err
			return
		}
		if currentWorker.Url == worker.Url {
			errChan <- errors.New("Worker already exists for this url.")
			return
		}
	}
	workerList = append(workerList, string(marshalledWorker))
	marshalledUpdatedWorker, err := json.Marshal(workerList)
	if err != nil {
		errChan <- err
		return
	}
	_, err = tx.Exec(context.Background(), "UPDATE workers SET workers=$1 WHERE app_name=$2", marshalledUpdatedWorker, appName)
	if err != nil {
		tx.Rollback(context.Background())
		errChan <- err
		return
	}
	errChan <- tx.Commit(context.Background())
}

func UpdateWorker(appName string, worker types.Worker) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	var workers string
	err = conn.QueryRow(context.Background(), "SELECT workers FROM workers WHERE app_name=$1", appName).Scan(&workers)
	if err != nil {
		return err
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		return err
	}
	for i, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			return err
		}
		if currentWorker.Url == worker.Url {
			tmp, _ := json.Marshal(worker)
			workerList[i] = string(tmp)
			break
		}
	}
	marshalledUpdatedWorker, err := json.Marshal(workerList)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(context.Background(), "UPDATE workers SET workers=$1 WHERE app_name=$2", marshalledUpdatedWorker, appName); err != nil {
		tx.Rollback(context.Background())
		return err
	}
	return tx.Commit(context.Background())
}

func DeleteWorker(appName string, url string) error {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return err
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	var workers string
	err = conn.QueryRow(context.Background(), "SELECT workers FROM workers WHERE app_name=$1", appName).Scan(&workers)
	if err != nil {
		return err
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		return err
	}
	for i, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			return err
		}
		if currentWorker.Url == url {
			workerList = slices.Delete(workerList, i, i+1)
			break
		}
	}
	marshalledUpdatedWorker, err := json.Marshal(workerList)
	if err != nil {
		return err
	}
	if _, err = tx.Exec(context.Background(), "UPDATE workers SET workers=$1 WHERE app_name=$2", marshalledUpdatedWorker, appName); err != nil {
		tx.Rollback(context.Background())
		return err
	}
	return tx.Commit(context.Background())
}

func GetWorkers(appName string) ([]types.Worker, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_USERS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(context.Background())
	var workers string
	err = conn.QueryRow(context.Background(), "SELECT workers FROM workers WHERE app_name=$1", appName).Scan(&workers)
	if err != nil {
		if err == pgx.ErrNoRows || workers == "" {
			return []types.Worker{}, nil
		}
		return nil, err
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		return nil, err
	}
	result := make([]types.Worker, 0, len(workerList))
	for _, workerEntry := range workerList {
		var worker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &worker); err != nil {
			return nil, err
		}
		result = append(result, worker)
	}
	return result, nil
}
