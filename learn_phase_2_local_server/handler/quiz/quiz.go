package quiz

import (
	"github.com/gin-gonic/gin"
)

// RegisterQuizRoutes registers all quiz routes to the given router group
func RegisterQuizRoutes(rg *gin.RouterGroup) {
	rg.GET("/", GetQuiz)
	rg.POST("/", PostQuiz)
	rg.PUT(":id", UpdateQuiz)
	rg.DELETE(":id", DeleteQuiz)
}
