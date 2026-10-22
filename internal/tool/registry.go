package tool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/openai/openai-go/v3"
)

type Registry struct {
	tools map[string]Tool
}

func NewRegistry() Registry {
	return Registry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool without replacing an existing registration.
// Register tools during startup, before the registry is used concurrently.
func (r *Registry) Register(t Tool) error {
	name := t.Definition.Name
	if strings.TrimSpace(name) == "" {
		return errors.New("register tool: name must not be empty")
	}
	if t.Execute == nil {
		return fmt.Errorf("register tool %q: execute handler must not be nil", name)
	}
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("register tool %q: already registered", name)
	}
	if r.tools == nil {
		r.tools = make(map[string]Tool)
	}
	r.tools[name] = t
	return nil
}

// Definitions returns API tool definitions sorted by name, ready for ChatCompletionNewParams.Tools.
func (r *Registry) Definitions() []openai.ChatCompletionToolUnionParam {
	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	sort.Strings(names)

	definitions := make([]openai.ChatCompletionToolUnionParam, 0, len(names))
	for _, name := range names {
		definitions = append(definitions, r.tools[name].ToParam())
	}
	return definitions
}

// Execute invokes a registered tool. The handler validates its arguments.
func (r *Registry) Execute(ctx context.Context, name string, args json.RawMessage) (string, error) {
	t, exists := r.tools[name]
	if !exists {
		return "", fmt.Errorf("execute tool %q: not registered", name)
	}
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("execute tool %q: %w", name, err)
	}
	result, err := t.Execute(ctx, args)
	if err != nil {
		return "", fmt.Errorf("execute tool %q: %w", name, err)
	}
	return result, nil
}
