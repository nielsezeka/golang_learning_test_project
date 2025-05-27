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
		// The /quiz group requires authentication and permission checks;
		// a valid token is required to access these routes.
		quizGroup := api.Group("/quiz")
		quizGroup.Use(auth.AuthMiddleware())
		quiz.RegisterQuizRoutes(quizGroup)
	}
	r.Static("/swagger_ui", "./swagger-ui/dist")
	return r
}
