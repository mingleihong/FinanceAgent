package service

import (
	"Finance/internal/core"
	"context"

	"github.com/openai/openai-go/v3"
)

type ChatService struct {
	agent *core.Agent
}

func NewChatService(agent *core.Agent) *ChatService {
	return &ChatService{
		agent: agent,
	}
}

func (chat *ChatService) Chat(ctx context.Context, prompt string, onDelta func(event string, content string) error) error {
	return chat.agent.Run(ctx, []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}, onDelta)
}
