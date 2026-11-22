package service

import (
	"Finance/internal/errs"
	"Finance/internal/model"
	"Finance/internal/repository"
	"context"
	"strconv"
)

type MessageService struct {
	messageRep *repository.MessageRep
}

func NewMessageService(messageRep *repository.MessageRep) *MessageService {
	return &MessageService{messageRep: messageRep}
}

// ListBySession 返回指定用户所属会话下的消息列表，按 id 升序。
// 归属过滤由 repository 层的 JOIN 条件完成，当前用户无法读取他人会话的消息。
func (s *MessageService) ListBySession(ctx context.Context, userID, sessionID int64) ([]model.Message, error) {
	entities, err := s.messageRep.ListBySessionIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}

	items := make([]model.Message, 0, len(entities))
	for _, entity := range entities {
		items = append(items, model.Message{
			ID:               strconv.FormatInt(entity.ID, 10),
			SessionID:        strconv.FormatInt(entity.SessionID, 10),
			Role:             entity.Role,
			Content:          entity.Content,
			ReasoningContent: entity.ReasoningContent,
			ToolCalls:        entity.ToolCalls,
			ToolCallID:       entity.ToolCallID,
			CreatedAt:        entity.CreatedAt,
		})
	}
	return items, nil
}
