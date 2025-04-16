package handler

import (
	"backend/internal/service"
	"context"
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
func GetColumns(c *gin.Context) {
    var req struct {
        Host     string `form:"host"`
        Port     string `form:"port"`
        Database string `form:"database"`
        User     string `form:"user"`
        JWTToken string `form:"jwt_token"`
        Table    string `form:"table"`
    }

    if err := c.Bind(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    conn, err := service.GetClickHouseClient(req.Host, req.Port, req.Database, req.User, req.JWTToken)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rows, err := conn.Query(context.Background(), "DESCRIBE TABLE "+req.Table)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    var columns []string
    for rows.Next() {
        var name, _type, _, _, _, _ string
        var skip1, skip2, skip3 string
        if err := rows.Scan(&name, &_type, &skip1, &skip2, &skip3); err == nil {
            columns = append(columns, name)
        }
    }

    c.JSON(http.StatusOK, gin.H{"columns": columns})
}

