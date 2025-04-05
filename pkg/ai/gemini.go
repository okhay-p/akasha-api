package ai

import (
	"akasha-api/pkg/config"
	"context"
	"log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var GeminiModel *genai.GenerativeModel
var Ctx context.Context

func InitGeminiModel(cfg config.Config) {
	Ctx = context.Background()
	apiKey := cfg.GeminiApiKey
	client, err := genai.NewClient(Ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}

	GeminiModel = client.GenerativeModel("gemini-2.0-flash")
}
