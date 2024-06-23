package main

import (
	"log"
	"net/http"

	"github.com/hanifmaliki/chat-app/internal/controller"
	"github.com/hanifmaliki/chat-app/internal/repository"
	routes "github.com/hanifmaliki/chat-app/internal/routes"
	"github.com/hanifmaliki/chat-app/internal/usecase"
	websocket "github.com/hanifmaliki/chat-app/internal/websocket"
	db "github.com/hanifmaliki/chat-app/pkg/db"
	util "github.com/hanifmaliki/chat-app/pkg/util"
)

func main() {
	util.LoadConfig()

	db, err := db.NewGormDB()
	if err != nil {
		log.Fatal("Error initializing storage: ", err)
	}

	messageRepo := repository.NewGormMessageRepository(db)
	chatUseCase := usecase.NewChatService(messageRepo)
	hub := websocket.NewHub()
	go hub.Run()

	chatController := controller.NewChatController(chatUseCase, hub)
	mux := http.NewServeMux()

	routes.SetupRoutes(mux, chatController)

	port := util.GetEnv("PORT", "8080")
	log.Printf("Server started on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
