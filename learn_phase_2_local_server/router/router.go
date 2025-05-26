package router

import (
	"learn_phase_2_local_server/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/login", handler.Login)
		api.POST("/refresh", handler.Refresh)
		// Protected routes
		api.Use(handler.AuthMiddleware())
		api.GET("/quiz", handler.GetQuiz)
		api.POST("/quiz", handler.PostQuiz)
		api.PUT("/quiz/:id", handler.UpdateQuiz)
		api.DELETE("/quiz/:id", handler.DeleteQuiz)
	}
	r.Static("/swagger_ui", "./swagger-ui/dist")
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, Gin!"})
	})

	return r
}
