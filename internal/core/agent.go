package core

import (
	"Finance/internal/tool"
	"context"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
)

const maxAgentSteps = 10

type Agent struct {
	LLM   *LLM
	Tools tool.Registry
}

func NewAgent(llm *LLM, tools tool.Registry) *Agent {
	return &Agent{
		LLM:   llm,
		Tools: tools,
	}
}

// Run streams reasoning and answer deltas, executing tools between model calls.
// It makes at most maxAgentSteps model calls and stops on the first error.
// Messages are copied so appending run history does not modify the caller's slice.
func (a *Agent) Run(ctx context.Context, messages []openai.ChatCompletionMessageParamUnion, onDelta func(event string, content string) error) error {
	history := append([]openai.ChatCompletionMessageParamUnion(nil), messages...)
	definitions := a.Tools.Definitions()

	for step := 0; step < maxAgentSteps; step++ {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("agent: %w", err)
		}

		message, err := a.LLM.StreamChat(ctx, history, definitions, onDelta)
		if err != nil {
			return fmt.Errorf("agent step %d: %w", step+1, err)
		}
		if len(message.ToolCalls) == 0 {
			return nil
		}

		// Do not execute tools if there is no model call left to consume their results.
		if step == maxAgentSteps-1 {
			break
		}
		history = append(history, message.ToParam())
		for _, call := range message.ToolCalls {
			if call.Type != "function" {
				return fmt.Errorf("agent: unsupported tool call type %q", call.Type)
			}
			result, err := a.Tools.Execute(ctx, call.Function.Name, json.RawMessage(call.Function.Arguments))
			if err != nil {
				return fmt.Errorf("agent tool call %q: %w", call.ID, err)
			}
			history = append(history, openai.ToolMessage(result, call.ID))
		}
	}

	return fmt.Errorf("agent: exceeded maximum of %d model calls", maxAgentSteps)
}
