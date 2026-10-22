package core

import (
	"context"
	"errors"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type LLM struct {
	client openai.Client
	model  string
}

func NewLLM(cfg *LLMConfig) *LLM {
	client := openai.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBaseURL(cfg.BaseURL),
	)
	return &LLM{
		client: client,
		model:  cfg.Model,
	}
}

// Chat returns the first choice's complete message, including any tool calls.
// Pass nil tools for a request without tool definitions.
func (l *LLM) Chat(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam) (*openai.ChatCompletionMessage, error) {
	completion, err := l.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    l.model,
		Messages: messages,
		Tools:    tools,
	})
	if err != nil {
		return nil, fmt.Errorf("chat: %w", err)
	}
	if len(completion.Choices) == 0 {
		return nil, errors.New("chat: response contains no choices")
	}
	return &completion.Choices[0].Message, nil
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
