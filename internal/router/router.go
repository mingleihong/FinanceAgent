package router

import (
	"Finance/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRouter(chatHandler *handler.ChatHandler, userHandler *handler.UserHandler) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/chat", chatHandler.Chat)
	}

	return r
}
