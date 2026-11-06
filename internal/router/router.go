package router

import (
	"Finance/internal/config"
	"Finance/internal/handler"
	"Finance/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewRouter(chatHandler *handler.ChatHandler, userHandler *handler.UserHandler, jwtConfig *config.JWTConfig) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/chat", middleware.JWTAuth(jwtConfig), chatHandler.Chat)
	}

	return r
}
