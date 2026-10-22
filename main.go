package main

import (
	"Finance/internal/core"
	"Finance/internal/handler"
	"Finance/internal/router"
	"Finance/internal/service"
)

func main() {
	cfg, err := core.Load()
	if err != nil {
		panic(err)
	}
	llm := core.NewLLM(cfg)
	chatService := service.NewChatService(llm)
	chatHandler := handler.NewChatHandler(chatService)

	r := router.NewRouter(chatHandler)
	if err = r.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
