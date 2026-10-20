package core

import (
	"Finance/internal/tool"
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type LLM struct {
	client openai.Client
	model  string
	tools  *tool.Registry
}

func NewLLM(cfg *LLMConfig, tools *tool.Registry) *LLM {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)
	return &LLM{
		client: client,
		model:  cfg.Model,
		tools:  tools,
	}
}

// Chat returns the first choice's complete text response.
func (l *LLM) Chat(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion) (string, error) {
	completion, err := l.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    l.model,
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("chat: %w", err)
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("chat: response contains no choices")
	}
	return completion.Choices[0].Message.Content, nil
}

// StreamChat synchronously delivers non-empty text deltas from the first choice.
// A callback error stops the stream. Text already delivered remains valid if an
// error occurs, but may be incomplete. Callers can cancel via ctx.
func (l *LLM) StreamChat(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion, onDelta func(string) error) error {
	stream := l.client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Model:    l.model,
			Messages: messages,
		},
	)
	defer stream.Close()

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) == 0 || chunk.Choices[0].Delta.Content == "" {
			continue
		}
		if err := onDelta(chunk.Choices[0].Delta.Content); err != nil {
			return fmt.Errorf("stream chat callback: %w", err)
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream chat: %w", err)
	}
	return nil
}
