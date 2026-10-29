package main

import (
	"Finance/internal/config"
	"Finance/internal/core"
	"Finance/internal/handler"
	"Finance/internal/router"
	"Finance/internal/service"
	"Finance/internal/tool"
	"fmt"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	llm := core.NewLLM(cfg.LLM)
	registry := tool.NewRegistry()
	registry.Register(tool.NewCurrentTimeTool())
	agent := core.NewAgent(llm, registry)
	chatService := service.NewChatService(agent)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler)
	if err = r.Run(fmt.Sprintf("%s:%d", cfg.App.Address, cfg.App.Port)); err != nil {
		panic(err)
	}
}
