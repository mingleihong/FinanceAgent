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
	tools := tool.NewRegistry()
	if err := tools.Register(tool.NewCurrentTimeTool()); err != nil {
		panic(err)
	}
	llm := core.NewLLM(cfg, tools)
	chatService := service.NewChatService(llm)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler)
	if err = r.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
