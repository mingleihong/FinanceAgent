package handler

import (
	"Finance/internal/errs"
	"Finance/internal/response"
	"Finance/internal/service"

	"github.com/gin-gonic/gin"
)

type chatRequest struct {
	Prompt string `json:"prompt" binding:"required"`
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

	err := chat.chatService.Chat(c.Request.Context(), req.Prompt, func(event string, content string) error {
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
