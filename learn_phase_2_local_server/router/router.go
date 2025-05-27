package router

import (
	"learn_phase_2_local_server/handler"
	"learn_phase_2_local_server/handler/auth"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", auth.Login)
		api.POST("/register", auth.Register)
		api.POST("/refresh", auth.Refresh)

		quizGroup := api.Group("/quiz")
		// quizGroup.Use(handler.AuthMiddleware())
		handler.RegisterQuizRoutes(quizGroup)
	}

	r.Static("/swagger_ui", "./swagger-ui/dist")
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})

	return r
}
