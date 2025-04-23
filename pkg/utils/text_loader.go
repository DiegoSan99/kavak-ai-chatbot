package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/load"
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/openai"
	"github.com/tmc/langchaingo/schema"
)

func LoadTextFileWithEmbedding(filePath string, openaiClient *openai.OpenAIClient) ([]schema.Document, error) {

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read text file: %w", err)
	}

	textContent := string(content)

	title := extractTitle(textContent)

	llmContainer := &load.LLMContainer{
		Embedder: openaiClient,
		EmbeddingConfig: load.EmbeddingConfig{
			ChunkSize:    1000,
			ChunkOverlap: 200,
		},
	}

	err = llmContainer.InitEmbedding()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize embedding model: %w", err)
	}

	textEmbedding := &load.LLMTextEmbedding{
		ChunkSize:    llmContainer.EmbeddingConfig.ChunkSize,
		ChunkOverlap: llmContainer.EmbeddingConfig.ChunkOverlap,
		Text:         textContent,
	}

	docs, keywords, inconsistentChunks, err := textEmbedding.SplitTextWithLLM()
	if err != nil {
		log.Printf("LLM text splitting failed: %v. Falling back to simple text splitting.", err)
		docs, err = textEmbedding.SplitText()
		if err != nil {
			return nil, fmt.Errorf("failed to split text: %w", err)
		}
	}

	for i, doc := range docs {
		doc.Metadata = map[string]any{
			"source": filePath,
			"title":  title,
			"chunk":  i + 1,
		}

		if len(keywords) > 0 {
			doc.Metadata["keywords"] = keywords
		}

		if inconsistentChunk, ok := inconsistentChunks[i]; ok {
			doc.Metadata["inconsistent"] = inconsistentChunk
		}

		docs[i] = doc
	}

	log.Printf("Loaded %d chunks from text file: %s using embedding functionality\n", len(docs), filePath)
	return docs, nil
}

func extractTitle(content string) string {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}

	return "Untitled"
}
