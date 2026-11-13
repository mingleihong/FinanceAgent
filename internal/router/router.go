package router

import (
	"Finance/internal/config"
	"Finance/internal/errs"
	"Finance/internal/handler"
	"Finance/internal/middleware"
	"Finance/internal/response"
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func NewRouter(chatHandler *handler.ChatHandler, userHandler *handler.UserHandler, jwtConfig *config.JWTConfig) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.CustomRecoveryWithWriter(nil, func(c *gin.Context, recovered any) {
		response.Fail(c, errs.Wrap(errs.ErrInternalServer, fmt.Errorf("panic: %v\n%s", recovered, debug.Stack())))
	}))
	r.HandleMethodNotAllowed = true
	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, errs.New(errs.ErrNotFound))
	})
	r.NoMethod(func(c *gin.Context) {
		response.Fail(c, errs.New(errs.ErrMethodNotAllowed))
	})
	api := r.Group("/api")
	{
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.POST("/chat", middleware.JWTAuth(jwtConfig), chatHandler.Chat)
	}

	return r
}
