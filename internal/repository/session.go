package repository

import (
	"Finance/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrSessionUserNotFound 表示创建会话时 user_id 在 users 表中不存在。
var ErrSessionUserNotFound = errors.New("session owner not found")

type SessionRep struct {
	db *pgxpool.Pool
}

func NewSessionRep(db *pgxpool.Pool) *SessionRep {
	return &SessionRep{db: db}
}

// Create 为指定用户创建会话，返回表实体，标题使用数据库默认值。
func (r *SessionRep) Create(ctx context.Context, userID int64) (*model.SessionEntity, error) {
	entity := &model.SessionEntity{}
	err := r.db.QueryRow(ctx, `
		INSERT INTO sessions (user_id)
		VALUES ($1)
		RETURNING id, user_id, title, created_at, updated_at
	`, userID).Scan(&entity.ID, &entity.UserID, &entity.Title, &entity.CreatedAt, &entity.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" && pgErr.ConstraintName == "fk_sessions_user" {
			return nil, ErrSessionUserNotFound
		}
		return nil, fmt.Errorf("create session: %w", err)
	}
	return entity, nil
}
