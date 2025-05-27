package handler

import "github.com/gin-gonic/gin"

// RegisterQuizRoutes registers quiz-related routes to the given router group.
func RegisterQuizRoutes(r *gin.RouterGroup) {
	r.GET("", GetQuiz)
	r.POST("", PostQuiz)
	r.PUT(":id", UpdateQuiz)
	r.DELETE(":id", DeleteQuiz)
}
