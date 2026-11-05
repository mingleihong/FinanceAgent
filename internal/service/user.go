package service

import (
	"Finance/internal/config"
	"Finance/internal/model"
	"Finance/internal/repository"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidRegistration = errors.New("invalid registration")
	ErrUsernameTaken       = errors.New("username already exists")
	ErrInvalidCredentials  = errors.New("invalid username or password")
)

type UserService struct {
	userRep   *repository.UserRep
	jwtConfig *config.JWTConfig
}

func NewUserService(userRep *repository.UserRep, jwtConfig *config.JWTConfig) *UserService {
	return &UserService{userRep: userRep, jwtConfig: jwtConfig}
}

func (s *UserService) Login(ctx context.Context, username, password string) (*model.LoginResult, error) {
	username = strings.TrimSpace(username)
	if !utf8.ValidString(username) || strings.ContainsRune(username, '\x00') || username == "" || utf8.RuneCountInString(username) > 64 || len(password) > 72 {
		return nil, ErrInvalidCredentials
	}
	user, err := s.userRep.FindByUsername(ctx, username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("verify password: %w", err)
	}
	now := time.Now().Truncate(time.Second)
	expiresAt := now.Add(s.jwtConfig.TTL).Truncate(time.Second)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.AccessTokenClaims{
		Username:  user.Username,
		Subject:   strconv.FormatInt(user.ID, 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})
	accessToken, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}
	return &model.LoginResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}, nil
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
