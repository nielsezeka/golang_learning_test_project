package utils

import (
	"github.com/gin-gonic/gin"
)

// ErrorReturnHandler centralizes error responses for handlers
func ErrorReturnHandler(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{"error": err.Error()})
}
