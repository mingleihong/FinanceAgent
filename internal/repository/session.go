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

// ListByUserID 分页查询用户的会话列表，按 (updated_at DESC, id DESC) 排序。
// offset 为跳过的行数，limit 为期望返回的最大行数。
func (r *SessionRep) ListByUserID(ctx context.Context, userID, offset, limit int64) ([]*model.SessionEntity, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY updated_at DESC, id DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}
	defer rows.Close()

	var entities []*model.SessionEntity
	for rows.Next() {
		entity := &model.SessionEntity{}
		if err := rows.Scan(&entity.ID, &entity.UserID, &entity.Title, &entity.CreatedAt, &entity.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan session: %w", err)
		}
		entities = append(entities, entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate sessions: %w", err)
	}
	return entities, nil
}

// CountByUserID 返回用户的会话总数，用于分页计算总页数。
func (r *SessionRep) CountByUserID(ctx context.Context, userID int64) (int64, error) {
	var total int64
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM sessions WHERE user_id = $1`, userID).Scan(&total); err != nil {
		return 0, fmt.Errorf("count sessions: %w", err)
	}
	return total, nil
}
