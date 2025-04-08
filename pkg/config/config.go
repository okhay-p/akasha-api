package config

import (
	"os"
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
