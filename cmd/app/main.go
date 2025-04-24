package main

import (
	"akasha-api/internal/handlers"
	"akasha-api/internal/middlewares"
	"akasha-api/internal/routes"
	"akasha-api/pkg/ai"
	"akasha-api/pkg/config"
	"akasha-api/pkg/db"
	"akasha-api/pkg/jwt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	db.InitDB(cfg)
	ai.InitGeminiModel(cfg)
	jwt.SetSecret(cfg)

	if !config.Dev {
		handlers.NewAuth()
	} else {
		log.Println("DEV MODE")
	}

	router := gin.Default()
	router.SetTrustedProxies([]string{})
	router.Use(middlewares.CORSMiddleware())
	routes.SetupRouter(router)
	routes.SetupFeedbackRoutes(router)

	router.Use(middlewares.AuthMiddleware())
	routes.SetupUserRoutes(router)
	routes.SetupTopicRoutes(router)
	routes.SetupQuestionRoutes(router)

	port := os.Getenv("PORT")
	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "8080"
	}
	if err := router.Run(host + ":" + port); err != nil {
		log.Panicf("error: %s", err)
	}

}
