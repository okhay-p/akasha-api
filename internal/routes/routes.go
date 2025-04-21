package routes

import (
	"akasha-api/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.GET("/", handlers.GetAPIStatus)
	router.POST("/login", handlers.HandleLogin)
	router.GET("/auth/:provider/callback", handlers.HandleOAuthCallback)
	router.GET("/auth/:provider", handlers.HandleOAuth)
	router.GET("/logout", handlers.HandleOAuthLogout)
}
