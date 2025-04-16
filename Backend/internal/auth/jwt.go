package auth

import (
	"backend/internal/config"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var configData = config.LoadConfig()

// AccessTokenClaims struct for the Access Token
type AccessTokenClaims struct {
	Sub   string `json:"sub"`   // User ID (subject)
	Role  string `json:"role"`  // User role
	jwt.RegisteredClaims
}

// RefreshTokenClaims struct for the Refresh Token
type RefreshTokenClaims struct {
	Sub   string `json:"sub"`  // User ID (subject)
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new Access Token (short-lived)
func GenerateAccessToken(userID string) (string, error) {
	claims := AccessTokenClaims{
		Sub: userID,
		Role: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(configData.AccessTokenExpire) * time.Hour)),
			Issuer:    "my-auth-service",
			Audience:  jwt.ClaimStrings{"clickhouse-proxy"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configData.JwtSecret))
}

// GenerateRefreshToken creates a new Refresh Token (long-lived)
func GenerateRefreshToken(userID string) (string, error) {
	claims := RefreshTokenClaims{
		Sub: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(configData.RefreshTokenExpire) * 24 * time.Hour)),
			Issuer:    "my-auth-service",
			Audience:  jwt.ClaimStrings{"clickhouse-proxy"},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(configData.JwtSecret))
}

// ValidateAccessToken validates an Access Token
func ValidateAccessToken(tokenStr string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configData.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken validates a Refresh Token
func ValidateRefreshToken(tokenStr string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(configData.JwtSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing refresh token: %v", err)
	}

	if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid refresh token")
}

func LoginAndSetTokens(c *gin.Context) {
	// Define the structure of the request body
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Get credentials from environment variables as fallback (if the user didn't send valid data)
	if credentials.Username == "" {
		credentials.Username = getEnv("CLICKHOUSE_USER", "admin")
	}
	if credentials.Password == "" {
		credentials.Password = getEnv("CLICKHOUSE_PASSWORD", "supersecret")
	}

	// Authenticate the user (this should be a real authentication mechanism)
	// In this example, we're assuming the username/password is "admin" and "supersecret"
	if credentials.Username != "admin" || credentials.Password != "supersecret" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate Access Token and Refresh Token for the user
	accessToken, err := GenerateAccessToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	refreshToken, err := GenerateRefreshToken(credentials.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating refresh token"})
		return
	}

	// Set the tokens as HTTP-only cookies
	SetTokensAsCookies(c, accessToken, refreshToken)

	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

// SetTokensAsCookies sets the Access and Refresh tokens as HTTP-only cookies
func SetTokensAsCookies(c *gin.Context, accessToken, refreshToken string) {
	// Set Access Token as a cookie (HTTP-only, Secure, SameSite=Strict)
	c.SetCookie("access_token", accessToken, 3600, "/", "", true, true )  // 1 hour expiry

	// Set Refresh Token as a cookie (HTTP-only, Secure, SameSite=Strict)
	c.SetCookie("refresh_token", refreshToken, 30*24*60*60, "/", "", true, true) // 30 days expiry
}
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}


// Refresh token handler: Issues a new access token using the refresh token from the cookie
func RefreshTokenHandler(c *gin.Context) {
	// Retrieve the refresh token from the cookie
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	// Validate the refresh token
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Invalid refresh token: %s", err.Error())})
		return
	}

	// Generate a new access token using the refresh token's user ID
	newAccessToken, err := GenerateAccessToken(claims.Sub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating access token"})
		return
	}

	// Set the new access token as a cookie (HTTP-only)
	c.SetCookie("access_token", newAccessToken, 3600, "/", "", false, true) // 1 hour expiry

	// Return the new access token
	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
