package utils

import (
    "encoding/csv"
    "io"
    "os"
    "time"
)

func ReadCSV(r io.Reader) ([]string, [][]any, error) {
    reader := csv.NewReader(r)
    reader.TrimLeadingSpace = true

    header, err := reader.Read()
    if err != nil {
        return nil, nil, err
    }

    var rows [][]any
    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, nil, err
        }

        var row []any
        for _, val := range record {
            row = append(row, val)
        }

        rows = append(rows, row)
    }

    return header, rows, nil
}

func WriteCSV(header []string, data [][]string) (string, error) {
    filePath := "output_" + time.Now().Format("150405") + ".csv"
    f, err := os.Create(filePath)
    if err != nil {
        return "", err
    }
    defer f.Close()

    writer := csv.NewWriter(f)
    writer.Write(header)
    writer.WriteAll(data)
    writer.Flush()
    return filePath, writer.Error()
}
