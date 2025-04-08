package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupQuestionRoutes(router *gin.Engine) {
	questionGroup := router.Group("/question")
	{
		questionGroup.GET("/:id/answer/:option", handlers.CheckAnswer)
	}

}
