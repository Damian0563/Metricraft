package enter

import (
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
)

func Connect() error {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", "localhost", 9000)},
		Auth: clickhouse.Auth{
			Database: "main",
			Username: "default",
		},
	})
	return conn.Ping()
}
