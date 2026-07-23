package db

import (
	"context"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "metricraft/proto/metricraft/proto"
	"time"
)

func InsertWorkerLog(ctx context.Context, url string, success bool) error {
	conn, err := getLogsPool()
	if err != nil {
		return err
	}
	status := 0
	if success {
		status = 1
	}
	_, err = conn.Exec(ctx, "INSERT INTO worker_logs (date, url, up) VALUES ($1, $2, $3)", time.Now(), url, status)
	return err
}

func DeleteWorkerlogs(ctx context.Context, url string) error {
	conn, err := getLogsPool()
	if err != nil {
		return err
	}
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
	conn, err := getLogsPool()
	if err != nil {
		return nil, err
	}
	loc := loadLocation(timezone)
	uptime := make([]*pb.WorkerUptimeEntry, 0)
	endDate := time.Now().In(loc).UTC().Add(-time.Hour * 24 * 31)
	res, err := conn.Query(ctx, "SELECT date, up FROM worker_logs WHERE url = $1 AND date > $2 ORDER BY date ASC", url, endDate)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	pollInterval := time.Minute * time.Duration(pollIntervalSetting)
	var lastPoll time.Time
	for res.Next() {
		var date time.Time
		var status int
		if err := res.Scan(&date, &status); err != nil {
			return nil, err
		}
		trimmedDate := date.Truncate(time.Minute)
		if !lastPoll.IsZero() {
			trimmedLastPoll := lastPoll.Truncate(time.Minute).Add(pollInterval)
			// Extra padding avoids treating near-duplicate polls (routing/timeouts) as missed intervals; min poll is 5m.
			for trimmedDate.After(trimmedLastPoll.Add(2 * time.Minute)) {
				if len(uptime) == 0 || !uptime[len(uptime)-1].Stamp.AsTime().Truncate(time.Minute).Equal(trimmedLastPoll) {
					uptime = append(uptime, &pb.WorkerUptimeEntry{
						Status: -1,
						Stamp:  timestamppb.New(trimmedLastPoll),
					})
				}
				trimmedLastPoll = trimmedLastPoll.Add(pollInterval)
			}
		}
		if len(uptime) > 0 && uptime[len(uptime)-1].Stamp.AsTime().Truncate(time.Minute).Equal(trimmedDate) {
			uptime[len(uptime)-1].Status = int32(status)
			uptime[len(uptime)-1].Stamp = timestamppb.New(trimmedDate)
		} else {
			uptime = append(uptime, &pb.WorkerUptimeEntry{
				Status: int32(status),
				Stamp:  timestamppb.New(trimmedDate),
			})
		}
		lastPoll = date
	}
	now := time.Now().In(loc).UTC()
	for !lastPoll.IsZero() && now.After(lastPoll.Add(pollInterval)) {
		lastPoll = lastPoll.Add(pollInterval)
		trimmed := lastPoll.Truncate(time.Minute)
		if len(uptime) > 0 && uptime[len(uptime)-1].Stamp.AsTime().Truncate(time.Minute).Equal(trimmed) {
			continue
		}
		uptime = append(uptime, &pb.WorkerUptimeEntry{
			Status: -1,
			Stamp:  timestamppb.New(trimmed),
		})
	}
	//delete old logs, do not throw error if occured
	_, _ = conn.Exec(ctx, "DELETE FROM worker_logs WHERE url = $1 AND date < $2", url, endDate)
	return &pb.WorkerUptime{Entries: uptime}, nil
}
