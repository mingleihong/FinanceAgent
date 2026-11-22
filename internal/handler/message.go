package handler

import (
	"Finance/internal/errs"
	"Finance/internal/response"
	"Finance/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// ListBySession 返回当前登录用户指定会话下的消息列表。
func (h *MessageHandler) ListBySession(c *gin.Context) {
	sessionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || sessionID <= 0 {
		response.Fail(c, errs.New(errs.ErrInvalidParam))
		return
	}
	userID := c.GetInt64("user_id")

	messages, err := h.messageService.ListBySession(c.Request.Context(), userID, sessionID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, messages)
}
