package main

import (
	"Finance/internal/config"
	"Finance/internal/core"
	"Finance/internal/db"
	"Finance/internal/handler"
	"Finance/internal/repository"
	"Finance/internal/router"
	"Finance/internal/service"
	"Finance/internal/tool"
	"context"
	"fmt"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := db.NewPostgresPool(ctx, cfg.Postgres)
	cancel()
	if err != nil {
		panic(err)
	}
	defer pool.Close()
	userRep := repository.NewUserRep(pool)
	userService := service.NewUserService(userRep, cfg.JWT)
	userHandler := handler.NewUserHandler(userService)

	llm := core.NewLLM(cfg.LLM)
	registry := tool.NewRegistry()
	registry.Register(tool.NewCurrentTimeTool())
	agent := core.NewAgent(llm, registry)
	chatService := service.NewChatService(agent)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler, userHandler)
	if err = r.Run(fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)); err != nil {
		panic(err)
	}
}
