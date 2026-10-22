package handler

import (
	"Finance/internal/service"
	"log"
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
			"code":    1,
			"message": "wrong request parameters",
		})
		return
	}

	data, err := chat.chatService.Chat(c.Request.Context(), req.Prompt)
	if err != nil {
		log.Printf("Error in chat service: %v\n", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code":    1,
			"message": "wrong service",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}
