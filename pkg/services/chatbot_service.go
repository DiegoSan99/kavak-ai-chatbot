package services

import (
	"context"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/openai"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/vectordb"
	"go.uber.org/zap"
)

type ChatbotService struct {
	vectorDb     *vectordb.RedisVectorDB
	logger       *zap.Logger
	openaiClient *openai.OpenAIClient
	chain        *ChatbotChain
}

func NewChatbotService(vectorDb *vectordb.RedisVectorDB, logger *zap.Logger, openaiClient *openai.OpenAIClient) *ChatbotService {
	llmClient, err := openaiClient.NewLLMClient()
	if err != nil {
		logger.Fatal("Failed to create OpenAI model", zap.Error(err))
	}

	chain := NewChatbotChain(llmClient, vectorDb, logger)

	return &ChatbotService{
		vectorDb:     vectorDb,
		logger:       logger,
		openaiClient: openaiClient,
		chain:        chain,
	}
}

func (s *ChatbotService) QueryChatbot(query string) (string, error) {
	s.logger.Info("Getting chatbot response", zap.String("query", query))

	ctx := context.Background()
	response, err := s.chain.Run(ctx, query)
	if err != nil {
		s.logger.Error("Error getting response from chain", zap.Error(err))
		return "", err
	}

	s.logger.Info("Successfully generated response for query")
	return response, nil
}
