package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hanifmaliki/chat-app/internal/usecase"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"
	"github.com/hanifmaliki/chat-app/pkg/websocket"

	gorilla_websocket "github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

var upgrader = gorilla_websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type ChatController struct {
	userUsecase *usecase.UserUseCase
	roomUsecase *usecase.RoomUseCase
	hub         *websocket.Hub
}

func NewChatController(userUsecase *usecase.UserUseCase, roomUsecase *usecase.RoomUseCase, hub *websocket.Hub) *ChatController {
	return &ChatController{
		userUsecase: userUsecase,
		roomUsecase: roomUsecase,
		hub:         hub,
	}
}

func (cc *ChatController) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to upgrade connection")
		response := pkg_model.Response[any]{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to upgrade connection",
		}
		json.NewEncoder(w).Encode(response)
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
		} else if strings.HasPrefix(string(message), "dm:") {
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
