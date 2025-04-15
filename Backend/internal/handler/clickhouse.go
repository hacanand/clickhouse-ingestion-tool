package handler

import (
    "backend/internal/service"
    "net/http"

    "github.com/gin-gonic/gin"
)

func ClickHouseToFile(c *gin.Context) {
    var req struct {
        Host      string   `form:"host"`
        Port      string   `form:"port"`
        Database  string   `form:"database"`
        User      string   `form:"user"`
        JWTToken  string   `form:"jwt_token"`
        Table     string   `form:"table"`
        Columns   []string `form:"columns[]"`
    }

    if err := c.Bind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    outputPath, err := service.ExportClickHouseToCSV(req.Host, req.Port, req.Database, req.User, req.JWTToken, req.Table, req.Columns)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Success", "file": outputPath})
}
