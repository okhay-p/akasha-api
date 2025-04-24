package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupTopicRoutes(router *gin.Engine) {
	topicGroup := router.Group("/topic")
	{
		topicGroup.GET("", handlers.GetTopicsL1)
		topicGroup.POST("", handlers.CreateTopic)
		topicGroup.GET("/details/:id", handlers.GetFullTopicDetails)
		topicGroup.GET("/progress/:id", handlers.FirstOrCreateTopicProgress)
		topicGroup.DELETE("/progress/:id", handlers.DeleteProgress)
		topicGroup.PUT("/progress/:id/:order", handlers.UpdateTopicProgressCurrentLesson)
		topicGroup.GET("/progress", handlers.GetTopicsRelatedToUser)
		topicGroup.GET("/:id", handlers.GetTopicByUUID)
	}
}
