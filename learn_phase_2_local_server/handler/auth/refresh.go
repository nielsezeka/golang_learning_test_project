package auth

import (
	"learn_phase_2_local_server/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Refresh godoc
//
// @Summary      Refresh JWT token
// @Description  Get a new access token using a refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh_token  body  object  true  "Refresh token"
// @Success      200  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /api/refresh [post]
func Refresh(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.APIError{Error: "Invalid request"})
		return
	}
	// Validate refresh token
	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, utils.APIError{Error: "Invalid refresh token"})
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		c.JSON(http.StatusUnauthorized, utils.APIError{Error: "Invalid refresh token claims"})
		return
	}
	userID := int(claims["user_id"].(float64))
	if refreshTokens[req.RefreshToken] != userID {
		c.JSON(http.StatusUnauthorized, utils.APIError{Error: "Refresh token not recognized"})
		return
	}
	newTokenString, err := createToken(userID, []byte(jwtSecret), time.Minute*15)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.APIError{Error: "Could not generate token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": newTokenString})
}
