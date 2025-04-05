package main

import (
	"akasha-api/pkg/config"
	"akasha-api/pkg/db"

	"gorm.io/gen"
)

func main() {
	// Connect to your database
	cfg := config.LoadConfig()
	db.InitDB(cfg)

	// Create a new generator
	g := gen.NewGenerator(gen.Config{
		OutPath: "./internal/models",  // Output path for generated models
		Mode:    gen.WithDefaultQuery, // Generate default query methods
	})

	// Use the database connection
	g.UseDB(db.DB)

	// Generate models based on existing tables
	g.GenerateAllTable()

	// Save the generated code
	g.Execute()
}
