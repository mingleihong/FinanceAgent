package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/openai/openai-go/v3"
)

// NewCurrentTimeTool creates a tool that returns the server's current time.
func NewCurrentTimeTool() Tool {
	return Tool{
		Definition: openai.FunctionDefinitionParam{
			Name:        "get_current_time",
			Description: openai.String("Get the server's current date and time in RFC3339 format, including the UTC offset."),
			Parameters: openai.FunctionParameters{
				"type":                 "object",
				"properties":           map[string]any{},
				"required":             []string{},
				"additionalProperties": false,
			},
			Strict: openai.Bool(true),
		},
		Execute: getCurrentTime,
	}
}

func getCurrentTime(_ context.Context, args json.RawMessage) (string, error) {
	var params map[string]json.RawMessage
	if err := json.Unmarshal(args, &params); err != nil {
		return "", fmt.Errorf("decode current time arguments: %w", err)
	}
	if params == nil || len(params) != 0 {
		return "", errors.New("current time arguments must be an empty JSON object")
	}
	return time.Now().Format(time.RFC3339), nil
}
