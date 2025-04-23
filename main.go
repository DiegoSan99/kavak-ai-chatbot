package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/config"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/openai"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/services"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/utils"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/vectordb"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/web"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/schema"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load()

	redisURL := os.Getenv("REDIS_URL")
	redisHost := redisURL
	redisPassword := ""

	if strings.Contains(redisURL, "@") {
		parts := strings.Split(redisURL, "@")
		auth := strings.TrimPrefix(parts[0], "redis://")
		if strings.Contains(auth, ":") {
			authParts := strings.Split(auth, ":")
			redisPassword = authParts[1]
		}
		redisHost = parts[1]
	} else {
		redisHost = strings.TrimPrefix(redisURL, "redis://")
	}

	openaiClient := &openai.OpenAIClient{
		Config: *config.LoadConfig(),
	}

	_, err := openaiClient.NewLLMClient()
	if err != nil {
		log.Fatalf("Failed to create OpenAI model: %v", err)
	}

	loadData := os.Getenv("LOAD_DATA")
	redisVectorDB := vectordb.NewRedisVectorDB(openaiClient, redisHost, redisPassword)
	index := "kavak-chatbot-index"
	if loadData == "true" {
		var allDocuments []schema.Document

		csvDocuments, err := utils.LoadCSVToDocuments("sample_caso_ai_engineer.csv")
		if err != nil {
			log.Printf("Failed to load CSV data: %v", err)
		} else {
			allDocuments = append(allDocuments, csvDocuments...)
			fmt.Printf("Loaded %d car records from CSV\n", len(csvDocuments))
		}

		textDocuments, err := utils.LoadTextFileWithEmbedding("value_proposal.txt", openaiClient)
		if err != nil {
			log.Printf("Failed to load text data: %v", err)
		} else {
			allDocuments = append(allDocuments, textDocuments...)
			fmt.Printf("Loaded %d chunks from value proposal text using embedding\n", len(textDocuments))
		}

		if len(allDocuments) > 0 {
			docIDs, err := redisVectorDB.AddDocuments(index, allDocuments)
			if err != nil {
				log.Fatalf("Failed to add documents to Redis: %v", err)
			}

			fmt.Printf("Added %d documents to Redis with IDs: %v\n", len(docIDs), docIDs)
		} else {
			fmt.Println("No documents were loaded, skipping Redis upload")
		}
	} else {
		fmt.Println("Skipping data loading as LOAD_DATA is set to false")
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	chatbotService := services.NewChatbotService(redisVectorDB, logger, openaiClient)

	e := echo.New()
	web.NewChatbotController(e, sugar, chatbotService)

	if err := e.Start(":8080"); err != nil {
		sugar.Fatalf("Failed to start server: %v", err)
	}
}
