package repository

import (
	"errors"

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
