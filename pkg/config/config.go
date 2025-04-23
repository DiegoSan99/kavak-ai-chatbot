package config

import "os"

type LLMConfig struct {
	Apiurl   string // OpenAI API endpoint
	AiModel  string // OpenAI model name
	APIToken string // OpenAI API key
}

func LoadConfig() *LLMConfig {
	return &LLMConfig{
		Apiurl:   os.Getenv("OPENAI_API_URL"),
		AiModel:  os.Getenv("OPENAI_MODEL"),
		APIToken: os.Getenv("OPENAI_API_KEY"),
	}
}
