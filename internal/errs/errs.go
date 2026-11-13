package errs

import (
	"errors"
	"fmt"
	"net/http"
)

type Code int

const (
	Success Code = 0

	// 通用 10xxx
	ErrInternalServer   Code = 10000
	ErrInvalidParam     Code = 10001
	ErrUnauthorized     Code = 10002
	ErrNotFound         Code = 10003
	ErrMethodNotAllowed Code = 10004

	// 用户 20xxx
	ErrUserExists         Code = 20001
	ErrInvalidCredentials Code = 20002

	// 聊天 30xxx
	ErrChatFailed Code = 30001
)

// 每个错误码的默认文案 + HTTP 状态码，集中维护
type definition struct {
	msg      string
	httpCode int
}

var defs = map[Code]definition{
	ErrInternalServer:   {"服务器内部错误", http.StatusInternalServerError},
	ErrInvalidParam:     {"参数错误", http.StatusBadRequest},
	ErrUnauthorized:     {"未登录或登录已过期", http.StatusUnauthorized},
	ErrNotFound:         {"资源不存在", http.StatusNotFound},
	ErrMethodNotAllowed: {"请求方法不支持", http.StatusMethodNotAllowed},

	ErrUserExists:         {"用户已存在", http.StatusConflict},
	ErrInvalidCredentials: {"用户名或密码错误", http.StatusUnauthorized},

	ErrChatFailed: {"回答生成失败，请稍后重试", http.StatusInternalServerError},
}

// AppError 是贯穿全链路的业务错误
type AppError struct {
	Code     Code
	Msg      string // 对外消息，可安全返回前端
	HTTPCode int
	Err      error // 原始错误，仅用于日志，绝不返回前端
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d]%s: %v", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("[%d]%s", e.Code, e.Msg)
}

func (e *AppError) Unwrap() error { return e.Err }

// ---- 构造函数 ----

// New 使用默认文案
func New(code Code) *AppError {
	d, ok := defs[code]
	if !ok {
		d = defs[ErrInternalServer]
		code = ErrInternalServer
	}
	return &AppError{Code: code, Msg: d.msg, HTTPCode: d.httpCode}
}

// Newf 覆盖文案（如"用户 %s 不存在"）
func Newf(code Code, format string, args ...any) *AppError {
	e := New(code)
	e.Msg = fmt.Sprintf(format, args...)
	return e
}

// Wrap 携带底层原始错误（日志用）
func Wrap(code Code, err error) *AppError {
	e := New(code)
	e.Err = err
	return e
}

// From 提取错误链中的 AppError；未找到（包括 nil）时包装为内部错误。
func From(err error) *AppError {
	var e *AppError
	if errors.As(err, &e) && e != nil {
		return e
	}
	return Wrap(ErrInternalServer, err)
}
