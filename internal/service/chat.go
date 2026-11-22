package service

import (
	"Finance/internal/core"
	"Finance/internal/errs"
	"Finance/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/openai/openai-go/v3"
)

type ChatService struct {
	agent      *core.Agent
	sessionRep *repository.SessionRep
	messageRep *repository.MessageRep
}

func NewChatService(agent *core.Agent, sessionRep *repository.SessionRep, messageRep *repository.MessageRep) *ChatService {
	return &ChatService{
		agent:      agent,
		sessionRep: sessionRep,
		messageRep: messageRep,
	}
}

// Chat 校验会话归属，从数据库重建历史，执行 Agent，并把本次完整交互落库。
// 客户端提交的 messages 不参与上下文，避免伪造历史。
func (chat *ChatService) Chat(ctx context.Context, userID, sessionID int64, prompt string, onDelta func(event string, content string) error) error {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return errs.Newf(errs.ErrInvalidParam, "消息不能为空")
	}

	exists, err := chat.sessionRep.ExistsByIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return errs.Wrap(errs.ErrInternalServer, err)
	}
	if !exists {
		return errs.Newf(errs.ErrNotFound, "会话不存在")
	}

	msgEntity, err := chat.messageRep.ListBySessionIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return errs.Wrap(errs.ErrInternalServer, err)
	}
	history := make([]openai.ChatCompletionMessageParamUnion, 0, len(msgEntity)+1)
	for _, v := range msgEntity {
		switch v.Role {
		case "system":
			history = append(history, openai.SystemMessage(v.Content))
		case "user":
			history = append(history, openai.UserMessage(v.Content))
		case "assistant":
			assistant := openai.AssistantMessage(v.Content).OfAssistant
			if len(v.ToolCalls) > 0 {
				if err := json.Unmarshal(v.ToolCalls, &assistant.ToolCalls); err != nil {
					return errs.Wrap(errs.ErrInternalServer, fmt.Errorf("decode tool calls for message %d: %w", v.ID, err))
				}
			}
			history = append(history, openai.ChatCompletionMessageParamUnion{OfAssistant: assistant})
		case "tool":
			if v.ToolCallID == nil || *v.ToolCallID == "" {
				return errs.Newf(errs.ErrInternalServer, "tool message %d has no tool_call_id", v.ID)
			}
			history = append(history, openai.ToolMessage(v.Content, *v.ToolCallID))
		default:
			return errs.Newf(errs.ErrInternalServer, "message %d has unsupported role %q", v.ID, v.Role)
		}
	}
	history = append(history, openai.UserMessage(prompt))

	err = chat.agent.Run(ctx, history, onDelta)
	if err != nil {
		return errs.Wrap(errs.ErrChatFailed, err)
	}
	return nil
}
