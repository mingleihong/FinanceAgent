package core

import (
	"context"
	"encoding/json"
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

// StreamChat streams reasoning and answer deltas while accumulating a message for
// the agent's next tool round. Thinking text is read only from reasoning_content.
func (l *LLM) StreamChat(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion, tools []openai.ChatCompletionToolUnionParam, onDelta func(event string, content string) error) (*openai.ChatCompletionMessage, error) {
	stream := l.client.Chat.Completions.NewStreaming(
		ctx,
		openai.ChatCompletionNewParams{
			Model:    l.model,
			Messages: messages,
			Tools:    tools,
		},
	)
	defer stream.Close()
	var accumulator openai.ChatCompletionAccumulator

	for stream.Next() {
		chunk := stream.Current()
		if !accumulator.AddChunk(chunk) {
			return nil, errors.New("stream chat: invalid completion chunk")
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		choice := chunk.Choices[0]
		if field, ok := choice.Delta.JSON.ExtraFields["reasoning_content"]; ok {
			var text string
			if err := json.Unmarshal([]byte(field.Raw()), &text); err == nil && text != "" {
				if err := onDelta("reasoning_content", text); err != nil {
					return nil, fmt.Errorf("stream chat callback: %w", err)
				}
			}
		}
		text := choice.Delta.Content
		if text == "" {
			text = choice.Delta.Refusal
		}
		if text != "" {
			if err := onDelta("content", text); err != nil {
				return nil, fmt.Errorf("stream chat callback: %w", err)
			}
		}
	}

	if err := stream.Err(); err != nil {
		return nil, fmt.Errorf("stream chat: %w", err)
	}
	if len(accumulator.Choices) == 0 {
		return nil, errors.New("stream chat: response contains no choices")
	}
	choice := accumulator.Choices[0]
	if choice.FinishReason != "stop" && choice.FinishReason != "tool_calls" {
		return nil, fmt.Errorf("stream chat: incomplete response (finish reason %q)", choice.FinishReason)
	}
	return &choice.Message, nil
}
