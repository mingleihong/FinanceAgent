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

// List 分页返回当前用户的会话列表。
// page 从 1 开始，offset = (page-1)*page_size；total 为会话总数。
func (s *SessionService) List(ctx context.Context, userID, page, pageSize int64) (*model.SessionPage, error) {
	offset := (page - 1) * pageSize
	entities, err := s.sessionRep.ListByUserID(ctx, userID, offset, pageSize)
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}
	total, err := s.sessionRep.CountByUserID(ctx, userID)
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}

	items := make([]model.Session, 0, len(entities))
	for _, entity := range entities {
		items = append(items, model.Session{
			ID:        strconv.FormatInt(entity.ID, 10),
			Title:     entity.Title,
			CreatedAt: entity.CreatedAt,
			UpdatedAt: entity.UpdatedAt,
		})
	}
	return &model.SessionPage{List: items, Total: total, Page: page, PageSize: pageSize}, nil
}
