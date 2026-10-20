package core

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type LLMConfig struct {
	APIKey  string
	BaseURL string
	Model   string
}

func Load() (*LLMConfig, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("load .env: %w", err)
	}
	cfg := LLMConfig{
		APIKey:  os.Getenv("LLM_API_KEY"),
		BaseURL: os.Getenv("LLM_BASE_URL"),
		Model:   os.Getenv("LLM_MODEL"),
	}

	return &cfg, nil
}
