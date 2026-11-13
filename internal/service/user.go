package service

import (
	"Finance/internal/config"
	"Finance/internal/errs"
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

type UserService struct {
	userRep   *repository.UserRep
	jwtConfig *config.JWTConfig
}

func NewUserService(userRep *repository.UserRep, jwtConfig *config.JWTConfig) *UserService {
	return &UserService{userRep: userRep, jwtConfig: jwtConfig}
}

func (s *UserService) Login(ctx context.Context, username, password string) (*model.LoginResult, error) {
	username = strings.TrimSpace(username)
	if !utf8.ValidString(username) || strings.ContainsRune(username, '\x00') || username == "" || utf8.RuneCountInString(username) > 64 || password == "" || len(password) > 72 {
		return nil, errs.New(errs.ErrInvalidParam)
	}
	user, err := s.userRep.FindByUsername(ctx, username)
	if errors.Is(err, repository.ErrUserNotFound) {
		return nil, errs.Wrap(errs.ErrInvalidCredentials, err)
	}
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errs.Wrap(errs.ErrInvalidCredentials, err)
		}
		return nil, errs.Wrap(errs.ErrInternalServer, fmt.Errorf("verify password: %w", err))
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
		return nil, errs.Wrap(errs.ErrInternalServer, fmt.Errorf("sign access token: %w", err))
	}
	return &model.LoginResult{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	}, nil
}

func (s *UserService) Register(ctx context.Context, username, password string) (*model.User, error) {
	username = strings.TrimSpace(username)
	if !utf8.ValidString(username) || strings.ContainsRune(username, '\x00') || username == "" || utf8.RuneCountInString(username) > 64 {
		return nil, errs.Newf(errs.ErrInvalidParam, "用户名须为 1–64 个字符，且不能包含空字符")
	}
	if password == "" {
		return nil, errs.Newf(errs.ErrInvalidParam, "密码不能为空")
	}
	if len(password) > 72 {
		return nil, errs.Newf(errs.ErrInvalidParam, "密码不能超过 72 字节")
	}
	// GenerateFromPassword generates a random salt and embeds it in the hash.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, fmt.Errorf("hash password: %w", err))
	}
	user, err := s.userRep.Create(ctx, username, string(hash))
	if errors.Is(err, repository.ErrUsernameTaken) {
		return nil, errs.Wrap(errs.ErrUserExists, err)
	}
	if err != nil {
		return nil, errs.Wrap(errs.ErrInternalServer, err)
	}
	return user, nil
}
