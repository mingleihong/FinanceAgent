package repository

import (
	"Finance/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrMessageSessionNotFound 表示新增消息时 session_id 在 sessions 表中不存在。
var ErrMessageSessionNotFound = errors.New("message session not found")

type MessageRep struct {
	db *pgxpool.Pool
}

func NewMessageRep(db *pgxpool.Pool) *MessageRep {
	return &MessageRep{db: db}
}

// ListBySessionIDAndUserID 查询指定用户所属会话下的全部消息，按 id 升序返回。
// 通过 JOIN sessions 同时施加 session_id 与 user_id 条件，
// 用户只能读取自己会话的消息，越权请求得到空列表而非泄露他人数据。
func (r *MessageRep) ListBySessionIDAndUserID(ctx context.Context, sessionID, userID int64) ([]*model.MessageEntity, error) {
	rows, err := r.db.Query(ctx, `
		SELECT m.id, m.session_id, m.role, m.content, m.reasoning_content, m.tool_calls, m.tool_call_id, m.created_at
		FROM messages m
		JOIN sessions s ON s.id = m.session_id
		WHERE m.session_id = $1 AND s.user_id = $2
		ORDER BY m.id ASC
	`, sessionID, userID)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}
	defer rows.Close()

	var entities []*model.MessageEntity
	for rows.Next() {
		entity := &model.MessageEntity{}
		if err := rows.Scan(
			&entity.ID,
			&entity.SessionID,
			&entity.Role,
			&entity.Content,
			&entity.ReasoningContent,
			&entity.ToolCalls,
			&entity.ToolCallID,
			&entity.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate messages: %w", err)
	}
	return entities, nil
}
