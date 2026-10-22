package main

import (
	"Finance/internal/core"
	"Finance/internal/handler"
	"Finance/internal/router"
	"Finance/internal/service"
	"Finance/internal/tool"
)

func main() {
	cfg, err := core.Load()
	if err != nil {
		panic(err)
	}
	llm := core.NewLLM(cfg)
	registry := tool.NewRegistry()
	registry.Register(tool.NewCurrentTimeTool())
	agent := core.NewAgent(llm, registry)
	chatService := service.NewChatService(agent)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler)
	if err = r.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
