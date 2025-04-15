package main

import (
    "backend/internal/handler"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.MaxMultipartMemory = 8 << 20 // 8 MiB

    api := r.Group("/api")
    {
        api.POST("/clickhouse-to-file", handler.ClickHouseToFile)
        api.POST("/file-to-clickhouse", handler.FileToClickHouse)
    }

	r.Static("/", "./") // Add this line to serve CSV files from root
    r.Run(":8080")
}
