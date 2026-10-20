package handler

import (
	"Finance/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type chatRequest struct {
	Prompt string `json:"prompt"`
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
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"success": "wrong",
		})
		return
	}

	data, err := chat.chatService.Chat(c.Request.Context(), req.Prompt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"success": "wrong",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": 1,
		"data":    data,
	})
}
