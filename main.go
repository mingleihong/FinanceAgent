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
	"log"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load application config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	pool, err := db.NewPostgresPool(ctx, cfg.Postgres)
	cancel()
	if err != nil {
		return fmt.Errorf("initialize database: %w", err)
	}
	defer pool.Close()
	userRep := repository.NewUserRep(pool)
	userService := service.NewUserService(userRep, cfg.JWT)
	userHandler := handler.NewUserHandler(userService)

	llm := core.NewLLM(cfg.LLM)
	registry := tool.NewRegistry()
	if err := registry.Register(tool.NewCurrentTimeTool()); err != nil {
		return fmt.Errorf("initialize tools: %w", err)
	}
	agent := core.NewAgent(llm, registry)
	chatService := service.NewChatService(agent)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler, userHandler, cfg.JWT)
	if err = r.Run(fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)); err != nil {
		return fmt.Errorf("run HTTP server: %w", err)
	}
	return nil
}
