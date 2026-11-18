package service

import (
	"Finance/internal/errs"
	"Finance/internal/model"
	"Finance/internal/repository"
	"context"
	"errors"
	"strconv"
)

type SessionService struct {
	sessionRep *repository.SessionRep
}

func NewSessionService(sessionRep *repository.SessionRep) *SessionService {
	return &SessionService{sessionRep: sessionRep}
}

// Create 为当前登录用户创建一个新会话，标题使用默认值。
func (s *SessionService) Create(ctx context.Context, userID int64) (*model.Session, error) {
	entity, err := s.sessionRep.Create(ctx, userID)
	if errors.Is(err, repository.ErrSessionUserNotFound) {
		return nil, errs.Newf(errs.ErrNotFound, "用户不存在，无法创建会话")
	}
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}
	return &model.Session{
		ID:        strconv.FormatInt(entity.ID, 10),
		Title:     entity.Title,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
