package auth

import (
	"learn_phase_2_local_server/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the request body for user registration
// swagger:model
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RegisterResponse represents the response body for user registration
// swagger:model
type RegisterResponse struct {
	Message string `json:"message"`
}

// Register godoc
//
// @Summary      Register new user
// @Description  Creates a new user with a hashed password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      handler.RegisterRequest  true  "User registration info"
// @Success      201   {object}  handler.RegisterResponse
// @Failure      400   {object}  map[string]string
// @Failure      409   {object}  map[string]string
// @Router       /api/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password required"})
		return
	}
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)", req.Username).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password+hashPassKey), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	_, err = db.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", req.Username, string(hash))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, RegisterResponse{Message: "User registered successfully"})
}
