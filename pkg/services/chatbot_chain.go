package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
	"go.uber.org/zap"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/prompts"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/vectordb"
)

type ChatbotChain struct {
	llmClient   llms.Model
	redisClient *vectordb.RedisVectorDB
	logger      *zap.Logger
	memory      []llms.MessageContent
}

func NewChatbotChain(llmClient llms.Model, redisClient *vectordb.RedisVectorDB, logger *zap.Logger) *ChatbotChain {
	return &ChatbotChain{
		llmClient:   llmClient,
		redisClient: redisClient,
		logger:      logger,
		memory:      make([]llms.MessageContent, 0),
	}
}

func (c *ChatbotChain) Run(ctx context.Context, query string) (string, error) {
	docs, err := c.redisClient.CosineSimilarity("kavak-chatbot-index", query, 5, 0.5)
	if err != nil {
		return "", fmt.Errorf("error getting similar documents: %v", err)
	}

	context := c.formatContext(docs)

	history := c.formatHistory()

	prompt := prompts.GetChatbotPromptWithHistory(context, query, history)

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, prompt),
	}

	messages = append(messages, llms.TextParts(llms.ChatMessageTypeHuman, query))

	response, err := c.llmClient.GenerateContent(ctx, messages)
	if err != nil {
		return "", fmt.Errorf("error getting response from LLM: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from LLM")
	}

	c.updateMemory(query, response.Choices[0].Content)

	return response.Choices[0].Content, nil
}

func (c *ChatbotChain) formatContext(docs []schema.Document) string {
	var context strings.Builder
	for _, doc := range docs {
		context.WriteString(doc.PageContent)
		context.WriteString("\n")
	}
	return context.String()
}

func (c *ChatbotChain) formatHistory() string {
	var history strings.Builder
	for _, msg := range c.memory {
		switch msg.Role {
		case llms.ChatMessageTypeHuman:
			history.WriteString(fmt.Sprintf("Human: %s\n", msg.Parts[0]))
		case llms.ChatMessageTypeAI:
			history.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Parts[0]))
		}
	}
	return history.String()
}

func (c *ChatbotChain) updateMemory(humanMessage string, aiMessage string) error {
	c.memory = append(c.memory, llms.TextParts(llms.ChatMessageTypeHuman, humanMessage))
	c.memory = append(c.memory, llms.TextParts(llms.ChatMessageTypeAI, aiMessage))
	if len(c.memory) > 10 {
		c.memory = c.memory[len(c.memory)-10:]
	}

	return nil
}
