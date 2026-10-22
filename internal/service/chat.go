package service

import (
	"Finance/internal/core"
	"context"

	"github.com/openai/openai-go/v3"
)

type ChatService struct {
	llm *core.LLM
}

func NewChatService(llm *core.LLM) *ChatService {
	return &ChatService{
		llm: llm,
	}
}

func (chat *ChatService) Chat(ctx context.Context, prompt string) (string, error) {
	data, err := chat.llm.Chat(ctx, []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}, nil)
	if err != nil {
		return "", err
	}
	return data.Content, nil
}
