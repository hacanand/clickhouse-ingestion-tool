package config

import (
	
	"os"
)

// Config holds all the application configuration
type Config struct {
	JwtSecret            string
	ClickhouseURL        string
	ClickhouseUser       string
	ClickhousePassword   string
	AccessTokenExpire    int64  // Access token expiration time in hours
	RefreshTokenExpire   int64  // Refresh token expiration time in days
	Port                 string
}

// LoadConfig loads environment variables into the Config struct
func LoadConfig() Config {
	return Config{
		JwtSecret:            getEnv("JWT_SECRET", "your-very-secret-key"),
		ClickhouseURL:        getEnv("CLICKHOUSE_URL", "http://localhost:8123"),
		ClickhouseUser:       getEnv("CLICKHOUSE_USER", "admin"),
		ClickhousePassword:   getEnv("CLICKHOUSE_PASSWORD", "supersecret"),
		AccessTokenExpire:    1,  // 1 hour
		RefreshTokenExpire:   30, // 30 days
		Port:                 getEnv("PORT", "3000"),
	}
}

// getEnv retrieves environment variables with a fallback default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
