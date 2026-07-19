package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"os"
	"time"
)

func InsertWorkerLog(ctx context.Context, url string, success bool) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	status := 0
	if success {
		status = 1
	}
	_, err = conn.Exec(ctx, "INSERT INTO worker_logs (date, url, up) VALUES ($1, $2, $3)", time.Now(), url, status)
	return err
}

func DeleteWorkerlogs(ctx context.Context, url string) error {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, "DELETE FROM worker_logs WHERE url=$1", url)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}
	err = tx.Commit(ctx)
	return err
}

func GetWorkerUptime(ctx context.Context, url string, timezone string, pollIntervalSetting int32) (*pb.WorkerUptime, error) {
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)
	loc := loadLocation(timezone)
	uptime := make([]*pb.WorkerUptimeEntry, 0)
	endDate := time.Now().In(loc).UTC().Add(-time.Hour * 24 * 31)
	res, err := conn.Query(ctx, "SELECT date, up FROM worker_logs WHERE url = $1 AND date > $2 ORDER BY date ASC", url, endDate)
	if err != nil {
		return nil, err
	}
	pollInterval := time.Minute * time.Duration(pollIntervalSetting)
	var lastPoll time.Time
	for res.Next() {
		var date time.Time
		var status int
		if err := res.Scan(&date, &status); err != nil {
			return nil, err
		}
		if !lastPoll.IsZero() {
			for date.After(lastPoll.Add(pollInterval)) {
				lastPoll = lastPoll.Add(pollInterval)
				uptime = append(uptime, &pb.WorkerUptimeEntry{
					Status: -1,
					Stamp:  timestamppb.New(lastPoll),
				})
			}
		}
		uptime = append(uptime, &pb.WorkerUptimeEntry{
			Status: int32(status),
			Stamp:  timestamppb.New(date),
		})
		lastPoll = date
	}
	now := time.Now().In(loc).UTC()
	for !lastPoll.IsZero() && now.After(lastPoll.Add(pollInterval)) {
		lastPoll = lastPoll.Add(pollInterval)
		uptime = append(uptime, &pb.WorkerUptimeEntry{
			Status: -1,
			Stamp:  timestamppb.New(lastPoll),
		})
	}
	//delete old logs, do not throw error if occured
	_, _ = conn.Exec(ctx, "DELETE FROM worker_logs WHERE url = $1 AND date < $2", url, endDate)
	return &pb.WorkerUptime{Entries: uptime}, nil
}
