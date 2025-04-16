package main

import (
	"fmt"
	"log"
	"net/http"
 
	"backend/internal/auth"         // Importing the auth package
	"backend/internal/config"       // Importing the config package
	"backend/internal/handler"      // Importing the handler package
	"github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"   // Importing the CORS middleware
)

var configData = config.LoadConfig()

func main() {
	r := gin.Default()
    r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Allow the frontend origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"}, // Allow specific HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Allow specific headers
		AllowCredentials: true, // Allow credentials (cookies)
	}))


	r.MaxMultipartMemory = 8 << 20 // 8 MiB



    r.POST("/api/login", auth.LoginAndSetTokens)
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "OK"})
    })
    api := r.Group("/api")
    api.Use(jwtMiddleware())
	{
		api.POST("/clickhouse-to-file", handler.ClickHouseToFile)
		api.POST("/file-to-clickhouse", handler.FileToClickHouse)
		api.POST("/get-columns", handler.GetColumns)
		api.POST("/refresh-token", auth.RefreshTokenHandler) // Endpoint to refresh access tokens
	}

	// Serve CSV files
	// r.Static("/", "./")
	
	// Start the server
	log.Printf("🚀 Server running on port %s", configData.Port)
	r.Run(fmt.Sprintf(":%s", configData.Port))
}

// JWT Middleware to validate Access Token from cookies
func jwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the cookies
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Validate Access Token
		_, err = auth.ValidateAccessToken(tokenStr)
		if err != nil {
			// If token expired, send an expiration error response
			if err.Error() == "access token expired" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Next()
	}
}
