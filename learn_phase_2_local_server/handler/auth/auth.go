package auth

import (
	"learn_phase_2_local_server/utils"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load() // Loads .env file if present
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))
var refreshSecret = []byte(os.Getenv("REFRESH_SECRET"))
var hashPassKey = os.Getenv("HASH_PASS_KEY")
var refreshTokens = make(map[string]int)

// AuthMiddleware checks for a valid JWT token in the Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.APIError{Error: "Missing or invalid Authorization header"})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.APIError{Error: "Invalid or expired token"})
			return
		}
		c.Set("claims", token.Claims)
		c.Next()
	}
}

// createToken generates a JWT token with the given userID, secret, and expiration duration
func createToken(userID int, secret []byte, duration time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	})
	return token.SignedString(secret)
}

// RegisterQuizRoutes registers quiz-related routes to the given router group.
func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.POST("/refresh", Refresh)
}
