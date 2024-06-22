package websocket

import (
	"net/http"

	"github.com/hanifmaliki/chat-app/internal/controller"
)

func SetupRoutes(mux *http.ServeMux, chatController *controller.ChatController) {
	mux.HandleFunc("/ws", chatController.HandleWebSocket)
}
