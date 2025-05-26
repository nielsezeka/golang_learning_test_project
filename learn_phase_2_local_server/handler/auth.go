package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-secret-key")
var refreshSecret = []byte("your-refresh-secret-key")
var refreshTokens = make(map[string]int) // map refreshToken -> userID

// Dummy user for demonstration
var demoUser = struct {
	Username string
	Password string // In production, store hashed passwords!
	ID       int
}{
	Username: "admin",
	Password: "password123", // In production, use a hashed password
	ID:       1,
}

// Login godoc
//	@Summary		User login
//	@Description	Authenticates user and returns JWT token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			credentials	body		object	true	"User credentials"
//	@Success		200			{object}	map[string]string
//	@Failure		401			{object}	map[string]string
//	@Router			/api/login [post]
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Dummy check (replace with DB lookup and hash check in production)
	if req.Username != demoUser.Username || req.Password != demoUser.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT (access token)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": demoUser.ID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(), // short-lived access token
	})
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Generate refresh token (longer-lived)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": demoUser.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})
	refreshTokenString, err := refreshToken.SignedString(refreshSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
		return
	}
	refreshTokens[refreshTokenString] = demoUser.ID

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "refresh_token": refreshTokenString})
}

// Refresh godoc
//	@Summary		Refresh JWT token
//	@Description	Get a new access token using a refresh token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			refresh_token	body		object	true	"Refresh token"
//	@Success		200				{object}	map[string]string
//	@Failure		401				{object}	map[string]string
//	@Router			/api/refresh [post]
func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	// Validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return refreshSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token claims"})
		return
	}
	userID := int(claims["user_id"].(float64))
	// Optionally check if refresh token is in store
	if refreshTokens[req.RefreshToken] != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not recognized"})
		return
	}
	// Generate new access token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Minute * 15).Unix(),
	})
	newTokenString, err := newToken.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}
		c.Set("claims", token.Claims)
		c.Next()
	}
}
