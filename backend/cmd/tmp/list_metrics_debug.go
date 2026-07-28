package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"backend/db"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load(".env")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_LOGS"))
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	db.SetLogsPool(pool)

	metrics, err := db.ListMetrics(ctx)
	if err != nil {
		panic(err)
	}
	out, _ := json.MarshalIndent(metrics, "", "  ")
	fmt.Println(string(out))

	if len(metrics) > 0 {
		rows, _ := pool.Query(ctx, "SELECT date, pg_typeof(date) FROM custom_metrics LIMIT 1")
		for rows.Next() {
			var d time.Time
			var typ string
			_ = rows.Scan(&d, &typ)
			fmt.Printf("\nraw date: %v loc=%v typ=%s formatted=%q\n", d, d.Location(), typ, d.UTC().Format(time.RFC3339))
		}
		rows.Close()
	}
}
