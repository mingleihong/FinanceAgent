package handler

import (
	"Finance/internal/errs"
	"Finance/internal/response"
	"Finance/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type chatRequest struct {
	Prompt    string `json:"prompt" binding:"required"`
	SessionID string `json:"session_id" binding:"required"`
}

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (chat *ChatHandler) Chat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errs.Wrap(errs.ErrInvalidParam, err))
		return
	}
	sessionID, err := strconv.ParseInt(req.SessionID, 10, 64)
	if err != nil || sessionID <= 0 {
		response.Fail(c, errs.New(errs.ErrInvalidParam))
		return
	}
	userID := c.GetInt64("user_id")

	err = chat.chatService.Chat(c.Request.Context(), userID, sessionID, req.Prompt, func(event string, content string) error {
		return response.SSE(c, event, gin.H{"content": content})
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	if err := response.SSE(c, "done", gin.H{}); err != nil {
		response.Fail(c, err)
	}
}
