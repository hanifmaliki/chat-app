package main

import (
	"log"
	"net/http"

	"github.com/hanifmaliki/chat-app/internal/controller"
	"github.com/hanifmaliki/chat-app/internal/repository"
	"github.com/hanifmaliki/chat-app/internal/usecase"
	"github.com/hanifmaliki/chat-app/internal/websocket"
	"github.com/hanifmaliki/chat-app/pkg/db"
	"github.com/hanifmaliki/chat-app/pkg/util"
)

func main() {
	util.LoadConfig()

	db, err := db.NewGormDB()
	if err != nil {
		log.Fatal("Error initializing storage: ", err)
	}

	messageRepo := repository.NewGormMessageRepository(db)
	chatService := usecase.NewChatService(messageRepo)
	hub := websocket.NewHub(chatService)
	go hub.Run()

	chatController := controller.NewChatController(hub)
	mux := http.NewServeMux()

	websocket.SetupRoutes(mux, chatController)

	port := util.GetEnv("PORT", "8080")
	log.Printf("Server started on :%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
