package openai

import (
	"github.com/DiegoSan99/kavak-document-preprocessor/pkg/config"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

type OpenAIClient struct {
	Config    config.LLMConfig
	LLMClient *openai.LLM
}

func (oc *OpenAIClient) NewEmbedder() (embeddings.Embedder, error) {
	// Create a new embedder with the text-embedding-3-small model
	llm, err := openai.New(
		openai.WithToken(oc.Config.APIToken),
		openai.WithBaseURL(oc.Config.Apiurl),
		openai.WithModel("text-embedding-3-small"), // Using the latest embedding model
	)
	if err != nil {
		return nil, err
	}
	return embeddings.NewEmbedder(llm)
}

func (oc *OpenAIClient) NewLLMClient() (llms.Model, error) {
	var err error
	oc.LLMClient, err = openai.New(
		openai.WithToken(oc.Config.APIToken),
		openai.WithBaseURL(oc.Config.Apiurl),
		openai.WithModel("gpt-3.5-turbo"),
	)
	return oc.LLMClient, err
}

func (oc *OpenAIClient) Initialized() bool {
	return oc.LLMClient != nil
}

func (oc *OpenAIClient) GetConfig() config.LLMConfig {
	return oc.Config
}
