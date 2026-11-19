package errs

import (
	"errors"
	"fmt"
)

// Code 业务错误码，按模块分段：1xxxx 通用，2xxxx 用户，3xxxx 聊天。
type Code int

const (
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

// messages 只维护对外文案；HTTP 状态码属于传输层，由 response 包映射。
var messages = map[Code]string{
	ErrInternalServer:     "服务器内部错误",
	ErrInvalidParam:       "参数错误",
	ErrUnauthorized:       "未登录或登录已过期",
	ErrNotFound:           "资源不存在",
	ErrMethodNotAllowed:   "请求方法不支持",
	ErrUserExists:         "用户已存在",
	ErrInvalidCredentials: "用户名或密码错误",
	ErrChatFailed:         "回答生成失败，请稍后重试",
}

// AppError 贯穿全链路的业务错误。
// Msg 可安全返回前端；Err 仅用于日志，绝不进入响应。
type AppError struct {
	Code Code
	Msg  string
	Err  error // 原始错误，仅用于日志
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d]%s: %v", e.Code, e.Msg, e.Err)
	}
	return fmt.Sprintf("[%d]%s", e.Code, e.Msg)
}

func (e *AppError) Unwrap() error { return e.Err }

func message(code Code) string {
	msg, ok := messages[code]
	if !ok {
		panic(fmt.Sprintf("errs: code %d has no message defined", code))
	}
	return msg
}

// New 用默认文案构造业务错误。
func New(code Code) *AppError {
	return &AppError{Code: code, Msg: message(code)}
}

// Newf 覆盖默认文案。format 必须是可直接返回前端的静态文案，
// 禁止拼接原始错误、SQL 或堆栈等内部信息。
func Newf(code Code, format string, args ...any) *AppError {
	e := New(code)
	e.Msg = fmt.Sprintf(format, args...)
	return e
}

// Wrap 把 err 归一化为 AppError：
// - err 链中已有 AppError 时原样返回，保留其业务码，避免被外层覆盖；
// - err 为普通错误时用 code 分类，原错误挂到链上仅用于日志。
// 用 Wrap(ErrInternalServer, err) 可作为无明确业务码时的兜底归一化。
func Wrap(code Code, err error) *AppError {
	if err == nil {
		return New(code)
	}
	var inner *AppError
	if errors.As(err, &inner) && inner != nil {
		return inner
	}
	return &AppError{Code: code, Msg: message(code), Err: err}
}

// CodeOf 提取错误链中的业务码。
func CodeOf(err error) (Code, bool) {
	var e *AppError
	if errors.As(err, &e) && e != nil {
		return e.Code, true
	}
	return 0, false
}
