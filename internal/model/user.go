package model

import "time"

// UserEntity 对应 users 表，供 repository 扫描。ID 与数据库 BIGINT 一致，使用 int64。
type UserEntity struct {
	ID           int64
	Username     string
	PasswordHash string `json:"-"` // 防误序列化泄漏哈希
	CreatedAt    time.Time
}

// User 是返回给前端的用户 DTO。ID 对外使用字符串，避免前端 JavaScript 处理大整数时精度丢失。
type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}
