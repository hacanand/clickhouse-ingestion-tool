package service

import (
	"fmt"
    "time"
    "github.com/ClickHouse/clickhouse-go/v2"
)

func GetClickHouseClient(host, port, db, user, token string) (clickhouse.Conn, error) {
    addr := fmt.Sprintf("%s:%s", host, port)
    return clickhouse.Open(&clickhouse.Options{
        Addr: []string{addr},
        Auth: clickhouse.Auth{
            Database: db,
            Username: user,
            Password: token,
        },
        DialTimeout: 5 * time.Second,
    })
}
