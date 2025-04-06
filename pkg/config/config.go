package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var AkashaLearnUUID string
var Dev bool
var FrontendUrl string

type Config struct {
	DBConnectionString string
	GeminiApiKey       string
	JwtSecret          string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AkashaLearnUUID = getEnv("AKASHALEARN_UUID", "")
	Dev = getEnv("DEV", "") == "true"
	FrontendUrl = getEnv("FRONTEND_URL", "*")

	return Config{
		DBConnectionString: getEnv("DB_CONNECTION_STRING", ""),
		GeminiApiKey:       getEnv("GEMINI_API_KEY", ""),
		JwtSecret:          getEnv("JWT_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
