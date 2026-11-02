package service

import (
	"Finance/internal/model"
	"Finance/internal/repository"
	"context"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegistration = errors.New("invalid registration")
	ErrUsernameTaken       = errors.New("username already exists")
)

type UserService struct {
	userRep *repository.UserRep
}

func NewUserService(userRep *repository.UserRep) *UserService {
	return &UserService{userRep: userRep}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	if !utf8.ValidString(username) || strings.ContainsRune(username, '\x00') || username == "" || utf8.RuneCountInString(username) > 64 {
		return nil, fmt.Errorf("%w: username must contain 1 to 64 characters and no null characters", ErrInvalidRegistration)
	}
	// GenerateFromPassword generates a random salt and embeds it in the hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errors.Is(err, bcrypt.ErrPasswordTooLong) {
		return nil, fmt.Errorf("%w: password must not exceed 72 bytes", ErrInvalidRegistration)
	}
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}
	user, err := s.userRep.Create(ctx, username, string(hash))
	if errors.Is(err, repository.ErrUsernameTaken) {
		return nil, ErrUsernameTaken
	}
	return user, err
}
