package main

import (
	"akasha-api/internal/middlewares"
	"akasha-api/internal/routes"
	"akasha-api/pkg/config"
	"akasha-api/pkg/db"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.LoadConfig()
	db.InitDB(cfg)

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	routes.SetupRouter(router)
	routes.SetupUserRoutes(router)
	if err := router.Run("localhost:8080"); err != nil {
		log.Fatal(err)
	}
}
