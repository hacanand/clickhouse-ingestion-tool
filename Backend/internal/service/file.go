package service

import (
	"backend/internal/utils"
	"context"
	"fmt"
	"io"
	"strings"
)

func ImportCSVToClickHouse(r io.Reader, host, port, db, user, token, table string) (int, error) {
    cols, data, err := utils.ReadCSV(r)
    if err != nil {
        return 0, err
    }

    conn, err := GetClickHouseClient(host, port, db, user, token)
    if err != nil {
        return 0, err
    }

    batch, err := conn.PrepareBatch(context.Background(), "INSERT INTO "+table+" ("+strings.Join(cols, ",")+")")
    if err != nil {
        return 0, err
    }

    for _, row := range data {
        if err := batch.Append(row...); err != nil {
            return 0, err
        }
    }

    err = batch.Send()
    return len(data), err
}

func ExportClickHouseToCSV(host, port, db, user, token, table string, columns []string) (string, error) {
    conn, err := GetClickHouseClient(host, port, db, user, token)
    if err != nil {
        return "", err
    }

    cols := strings.Join(columns, ",")
    rows, err := conn.Query(context.Background(), fmt.Sprintf("SELECT %s FROM %s", cols, table))
    if err != nil {
        return "", err
    }

    var allRows [][]string
    for rows.Next() {
        columns := rows.Columns()
        values := make([]interface{}, len(columns))
        valuePtrs := make([]interface{}, len(columns))
        for i := range values {
            valuePtrs[i] = &values[i]
        }

        if err := rows.Scan(valuePtrs...); err != nil {
            return "", err
        }

        strRow := make([]string, len(values))
        for i, val := range values {
            strRow[i] = fmt.Sprint(val)
        }
        allRows = append(allRows, strRow)
    }

    return utils.WriteCSV(columns, allRows)
}
