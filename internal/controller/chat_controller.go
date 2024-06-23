package controller

import (
	"log"
	"net/http"
	"strings"

	"github.com/hanifmaliki/chat-app/internal/usecase"
	"github.com/hanifmaliki/chat-app/internal/websocket"

	gorilla_websocket "github.com/gorilla/websocket"
)

var upgrader = gorilla_websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	usecase *usecase.ChatUseCase
	hub     *websocket.Hub
}

func NewChatController(usecase *usecase.ChatUseCase, hub *websocket.Hub) *ChatController {
	return &ChatController{
		usecase: usecase,
		hub:     hub,
	}
}

func (cc *ChatController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	clientID := r.URL.Query().Get("id")
	client := &websocket.Client{Hub: cc.hub, Conn: conn, Send: make(chan []byte, 256), ID: clientID}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump(func(c *websocket.Client, message []byte) {
		if strings.HasPrefix(string(message), "dm:") {
			parts := strings.SplitN(string(message), ":", 3)
			if len(parts) == 3 {
				targetID := parts[1]
				msg := parts[2]
				if targetClient, ok := c.Hub.ClientsByID[targetID]; ok {
					targetClient.Send <- []byte("dm:" + c.ID + ":" + msg)
				}
			}
		}
	})
}
