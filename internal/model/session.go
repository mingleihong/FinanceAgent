package model

import "time"

// SessionEntity 对应 sessions 表，供 repository 层扫描。
// ID 与数据库 BIGINT 一致，使用 int64。
type SessionEntity struct {
	ID        int64
	UserID    int64
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session 是返回给前端的会话 DTO。
// ID 对外使用字符串，避免前端 JavaScript 处理大整数时精度丢失。
type Session struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
