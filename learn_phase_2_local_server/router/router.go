package router

import (
	"learn_phase_2_local_server/handler/auth"
	"learn_phase_2_local_server/handler/quiz"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		authGroup := api.Group("/auth")
		auth.RegisterAuthRoutes(authGroup)
		quizGroup := api.Group("/quiz")
		quiz.RegisterQuizRoutes(quizGroup)
	}
	r.Static("/swagger_ui", "./swagger-ui/dist")
	return r
}
