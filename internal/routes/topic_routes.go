package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupTopicRoutes(router *gin.Engine) {
	topicGroup := router.Group("/topic")
	{
		topicGroup.POST("", handlers.CreateTopic)
		topicGroup.GET("/:id", handlers.GetTopicByUUID)
	}
}
