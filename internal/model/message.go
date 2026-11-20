package model

import (
	"encoding/json"
	"time"
)

// MessageEntity 对应 messages 表，供 repository 层扫描。
// ID 与数据库 BIGINT 一致，使用 int64；ToolCalls 保存 JSONB 原始字节。
type MessageEntity struct {
	ID               int64
	SessionID        int64
	Role             string
	Content          string
	ReasoningContent *string
	ToolCalls        []byte
	ToolCallID       *string
	CreatedAt        time.Time
}

// Message 是返回给前端的消息 DTO。
// ID 与 SessionID 对外使用字符串，避免前端 JavaScript 处理大整数时精度丢失。
type Message struct {
	ID               string          `json:"id"`
	SessionID        string          `json:"session_id"`
	Role             string          `json:"role"`
	Content          string          `json:"content"`
	ReasoningContent *string         `json:"reasoning_content,omitempty"`
	ToolCalls        json.RawMessage `json:"tool_calls,omitempty"`
	ToolCallID       *string         `json:"tool_call_id,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
}
