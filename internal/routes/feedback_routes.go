package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupFeedbackRoutes(router *gin.Engine) {
	feedbackGroup := router.Group("/feedback")
	{
		feedbackGroup.POST("", handlers.CreateFeedback)
	}

}
