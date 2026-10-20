package router

import (
	"Finance/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(chatHandler *handler.ChatHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/chat", chatHandler.Chat)
	}

	return r
}
