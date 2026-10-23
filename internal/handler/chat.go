package handler

import (
	"Finance/internal/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

	if strings.TrimSpace(req.Prompt) == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    1,
			"message": "prompt must not be empty",
		},
		)
		return
	}
	c.Header("Content-Type", "text/event-stream; charset=utf-8")
	c.Header("Cache-Control", "no-cache")
	c.Header("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	c.Writer.Flush()
	send := func(event string, payload any) error {
		if err := c.Request.Context().Err(); err != nil {
			return err
		}
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(c.Writer, "event: %s\ndata: %s\n\n", event, data); err != nil {
			return err
		}
		c.Writer.Flush()
		return nil
	}
	err := chat.chatService.Chat(c.Request.Context(), req.Prompt, func(event string, content string) error {
		return send(event, gin.H{"content": content})
	})
	if err != nil {
		log.Printf("Error in chat stream: %v", err)
		_ = send("error", gin.H{"message": "回答生成失败，请稍后重试"})
		return
	}
	_ = send("done", gin.H{})
}
