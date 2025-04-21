package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupUserRoutes(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.GET("", handlers.GetUsers)
		userGroup.GET("/profile", handlers.GetUserByUUID)
	}
}
