package controller

import (
	"log"
	"net/http"

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
	hub *websocket.Hub
}

func NewChatController(hub *websocket.Hub) *ChatController {
	return &ChatController{hub: hub}
}

func (cc *ChatController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	client := &websocket.Client{Hub: cc.hub, Conn: conn, Send: make(chan []byte, 256)}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
