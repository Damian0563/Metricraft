package db

import (
	"backend/types"
	"context"
	"encoding/json"
	"errors"
	"github.com/jackc/pgx/v5"
	"slices"
)

func SaveWorker(appName string, worker types.Worker, errChan chan error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		errChan <- err
		return
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		errChan <- err
		return
	}
	var workers string
	err = tx.QueryRow(ctx, "SELECT workers FROM workers WHERE app_name=$1 FOR UPDATE", appName).Scan(&workers)
	if err != nil {
		tx.Rollback(ctx)
		errChan <- err
		return
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		tx.Rollback(ctx)
		errChan <- err
		return
	}
	if len(workerList) >= 5 {
		tx.Rollback(ctx)
		errChan <- errors.New("Worker limit reached.")
		return
	}
	var marshalledWorker []byte
	if marshalledWorker, err = json.Marshal(worker); err != nil {
		tx.Rollback(ctx)
		errChan <- err
		return
	}
	for _, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			tx.Rollback(ctx)
			errChan <- err
			return
		}
		if currentWorker.Url == worker.Url {
			tx.Rollback(ctx)
			errChan <- errors.New("Worker already exists for this url.")
			return
		}
	}
	workerList = append(workerList, string(marshalledWorker))
	errChan <- UpdateWorkers(ctx, &tx, appName, workerList, workers)
}

func UpdateWorker(appName string, worker types.Worker) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	var workers string
	err = tx.QueryRow(ctx, "SELECT workers FROM workers WHERE app_name=$1 FOR UPDATE", appName).Scan(&workers)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		tx.Rollback(ctx)
		return err
	}
	for i, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			tx.Rollback(ctx)
			return err
		}
		if currentWorker.Url == worker.Url {
			tmp, err := json.Marshal(worker)
			if err != nil {
				tx.Rollback(ctx)
				return err
			}
			workerList[i] = string(tmp)
			break
		}
	}
	return UpdateWorkers(ctx, &tx, appName, workerList, workers)
}

func DeleteWorker(appName string, url string) error {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return err
	}
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	var workers string
	err = tx.QueryRow(ctx, "SELECT workers FROM workers WHERE app_name=$1 FOR UPDATE", appName).Scan(&workers)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	var workerList []string
	if err = json.Unmarshal([]byte(workers), &workerList); err != nil {
		tx.Rollback(ctx)
		return err
	}
	found := false
	for i, workerEntry := range workerList {
		var currentWorker types.Worker
		if err = json.Unmarshal([]byte(workerEntry), &currentWorker); err != nil {
			tx.Rollback(ctx)
			return err
		}
		if currentWorker.Url == url {
			workerList = slices.Delete(workerList, i, i+1)
			found = true
			break
		}
	}
	if !found {
		tx.Rollback(ctx)
		return errors.New("worker not found")
	}
	return UpdateWorkers(ctx, &tx, appName, workerList, workers)
}

func GetWorkers(appName string) ([]types.Worker, error) {
	ctx := context.Background()
	conn, err := getUsersPool()
	if err != nil {
		return nil, err
	}
	var workers string
	err = conn.QueryRow(ctx, "SELECT workers FROM workers WHERE app_name=$1", appName).Scan(&workers)
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

func UpdateWorkers(ctx context.Context, tx *pgx.Tx, appName string, workerList []string, previousWorkers string) error {
	marshalledUpdatedWorker, err := json.Marshal(workerList)
	if err != nil {
		(*tx).Rollback(ctx)
		return err
	}
	tag, err := (*tx).Exec(ctx, "UPDATE workers SET workers=$1 WHERE app_name=$2 AND workers::text=$3", marshalledUpdatedWorker, appName, previousWorkers)
	if err != nil {
		(*tx).Rollback(ctx)
		return err
	}
	if tag.RowsAffected() == 0 {
		(*tx).Rollback(ctx)
		return errors.New("workers changed concurrently")
	}
	return (*tx).Commit(ctx)
}
