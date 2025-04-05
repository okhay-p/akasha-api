package db

import (
	"akasha-api/pkg/config"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(cfg config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DBConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	log.Println("Database connection established")
}
