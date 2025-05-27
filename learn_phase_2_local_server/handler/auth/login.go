package auth

import (
	"learn_phase_2_local_server/db"
	"learn_phase_2_local_server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// LoginResponse represents the response body for a successful login
// swagger:model
type LoginResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

// Login godoc
//
// @Summary      User login
// @Description  Authenticates user and returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  object  true  "User credentials"
// @Success      200  {object}  handler.LoginResponse
// @Failure      401  {object}  map[string]string
// @Router       /api/login [post]
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: "Invalid request"})
		return
	}

	var user struct {
		ID           int
		Username     string
		PasswordHash string
	}
	err := db.DB.QueryRow("SELECT id, username, password FROM users WHERE username = $1", req.Username).
		Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.APIError{Error: "User does not existed"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password+hashPassKey)) != nil {
		c.JSON(http.StatusUnauthorized, utils.APIError{Error: "Invalid password"})
		return
	}

	userID := user.ID
	tokenString, err := createToken(userID, []byte(jwtSecret), time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: "Could not generate token"})
		return
	}

	refreshTokenString, err := createToken(userID, []byte(refreshSecret), time.Hour*24*7)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: "Could not generate refresh token"})
		return
	}
	refreshTokens[refreshTokenString] = userID

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "refresh_token": refreshTokenString})
}
