package handler

import (
	"Finance/internal/response"
	"Finance/internal/service"

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
