package main

import (
	"akasha-api/internal/middlewares"
	"akasha-api/internal/routes"
	"akasha-api/pkg/ai"
	"akasha-api/pkg/config"
	"akasha-api/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	db.InitDB(cfg)
	ai.InitGeminiModel(cfg)

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	routes.SetupRouter(router)
	routes.SetupUserRoutes(router)
	routes.SetupTopicRoutes(router)

	if err := router.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
