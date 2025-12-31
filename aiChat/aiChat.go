package aiChat

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/googleai"
)

func test() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := context.Background()
	llmClient, err := googleai.New(ctx,
		googleai.WithAPIKey(os.Getenv("GEMINI_API_KEY")),
		googleai.WithDefaultModel(os.Getenv("GEMINI_MODEL_ID")),
	)
	if err != nil {
		log.Fatal(err)
	}
	prompt := "What is the capital of France?"
	completion, err := llms.GenerateFromSinglePrompt(context.Background(), llmClient, prompt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(completion)
}
