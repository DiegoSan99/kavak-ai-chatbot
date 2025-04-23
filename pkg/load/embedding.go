package load

import (
	"strings"

	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/openai"
	"github.com/tmc/langchaingo/embeddings"
	lc "github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/schema"
)

// This embedder is used for the value proposal since csv is easier to split due to data being ordered
type LLMTextEmbedding struct {
	ChunkSize         int
	ChunkOverlap      int
	Text              string
	EmbeddedDocuments []schema.Document
}

type EmbeddingClient interface {
	NewEmbedder() (embeddings.Embedder, error)
	Initialized() bool
}

func (llm *LLMContainer) InitEmbedding() error {
	openaiLLM, err := lc.New(
		lc.WithToken(llm.Embedder.(*openai.OpenAIClient).Config.APIToken),
		lc.WithModel(llm.Embedder.(*openai.OpenAIClient).Config.AiModel),
	)
	if err != nil {
		return err
	}
	llm.Embedder.(*openai.OpenAIClient).LLMClient = openaiLLM

	return nil
}

func (lt *LLMTextEmbedding) SplitText() ([]schema.Document, error) {
	paragraphs := strings.Split(lt.Text, "\n\n")
	var docs []schema.Document

	for _, p := range paragraphs {
		if strings.TrimSpace(p) != "" {
			docs = append(docs, schema.Document{
				PageContent: p,
			})
		}
	}

	return docs, nil
}

func (lt *LLMTextEmbedding) SplitTextWithLLM() ([]schema.Document, []string, map[int]string, error) {
	docs, err := lt.SplitText()
	if err != nil {
		return nil, nil, nil, err
	}

	return docs, []string{}, map[int]string{}, nil
}
