package handler

import (
	"Finance/internal/errs"
	"Finance/internal/response"
	"Finance/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *service.SessionService
}

func NewSessionHandler(sessionService *service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// Create 为当前登录用户创建会话。
func (h *SessionHandler) Create(c *gin.Context) {
	userID := c.GetInt64("user_id")

	session, err := h.sessionService.Create(c.Request.Context(), userID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, session)
}

// List 分页返回当前登录用户的会话列表。
func (h *SessionHandler) List(c *gin.Context) {
	page, err := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	if err != nil || page < 1 || page > 100000 {
		response.Fail(c, errs.New(errs.ErrInvalidParam))
		return
	}
	pageSize, err := strconv.ParseInt(c.DefaultQuery("page_size", "10"), 10, 64)
	if err != nil || pageSize <= 0 || pageSize > 100 {
		response.Fail(c, errs.New(errs.ErrInvalidParam))
		return
	}
	userID := c.GetInt64("user_id")

	pageData, err := h.sessionService.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, pageData)
}
