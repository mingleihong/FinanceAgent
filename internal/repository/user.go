package repository

import (
	"Finance/internal/model"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUsernameTaken = errors.New("username already exists")
	ErrUserNotFound  = errors.New("user not found")
)

type UserRep struct {
	db *pgxpool.Pool
}

func NewUserRep(db *pgxpool.Pool) *UserRep {
	return &UserRep{db: db}
}

// FindByUsername returns the user including the hash for password verification.
func (r *UserRep) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(ctx, `
		SELECT id, username, created_at, password_hash
		FROM users
		WHERE username = $1
	`, username).Scan(&user.ID, &user.Username, &user.CreatedAt, &user.PasswordHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by username: %w", err)
	}
	return user, nil
}

func (r *UserRep) Create(ctx context.Context, username, passwordHash string) (*model.User, error) {
	user := &model.User{PasswordHash: passwordHash}
	err := r.db.QueryRow(ctx, `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, username, created_at
	`, username, passwordHash).Scan(&user.ID, &user.Username, &user.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "uk_users_username" {
			return nil, ErrUsernameTaken
		}
		return nil, fmt.Errorf("create user: %w", err)
	}
	return user, nil
}
