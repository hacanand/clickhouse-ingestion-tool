package handler

import (
    "backend/internal/service"
    "net/http"

    "github.com/gin-gonic/gin"
)

func FileToClickHouse(c *gin.Context) {
    host := c.PostForm("host")
    port := c.PostForm("port")
    database := c.PostForm("database")
    user := c.PostForm("user")
    jwtToken := c.PostForm("jwt_token")
    targetTable := c.PostForm("target_table")

    file, err := c.FormFile("file")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
        return
    }

    f, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer f.Close()

    count, err := service.ImportCSVToClickHouse(f, host, port, database, user, jwtToken, targetTable)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Ingested", "rows": count})
}
