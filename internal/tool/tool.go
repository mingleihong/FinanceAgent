package tool

import (
	"context"
	"encoding/json"

	"github.com/openai/openai-go/v3"
)

type Handler func(
	ctx context.Context,
	args json.RawMessage,
) (string, error)

type Tool struct {
	Definition openai.FunctionDefinitionParam
	Execute    Handler
}

// ToParam returns the API tool definition without the local execution handler.
func (t Tool) ToParam() openai.ChatCompletionToolUnionParam {
	return openai.ChatCompletionFunctionTool(t.Definition)
}
