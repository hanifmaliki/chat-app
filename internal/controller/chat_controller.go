package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hanifmaliki/chat-app/internal/entity"
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
	userUsecase     *usecase.UserUseCase
	roomUsecase     *usecase.RoomUseCase
	roomUserUsecase *usecase.RoomUserUseCase
	hub             *websocket.Hub
}

func NewChatController(userUsecase *usecase.UserUseCase, roomUsecase *usecase.RoomUseCase, roomUserUsecase *usecase.RoomUserUseCase, hub *websocket.Hub) *ChatController {
	return &ChatController{
		userUsecase:     userUsecase,
		roomUsecase:     roomUsecase,
		roomUserUsecase: roomUserUsecase,
		hub:             hub,
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
	username := r.URL.Query().Get("username")
	client := &websocket.Client{Hub: cc.hub, Conn: conn, Send: make(chan []byte, 256), Username: username}
	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump(func(c *websocket.Client, message []byte) {
		msg := string(message)
		if strings.HasPrefix(msg, "dm:") {
			parts := strings.SplitN(msg, ":", 3)
			if len(parts) == 3 {
				targetID := parts[1]
				dmMessage := parts[2]
				if targetClient, ok := c.Hub.ClientsByUsername[targetID]; ok {
					targetClient.Send <- []byte("DM from " + c.Username + ": " + dmMessage)
					c.Send <- []byte("Successfully sent DM to " + targetID)
				}
			}
		} else if strings.HasPrefix(msg, "room:create:") {
			parts := strings.SplitN(msg, ":", 3)
			if len(parts) == 3 {
				roomName := parts[2]
				err := cc.roomUsecase.Create(&entity.Room{Name: roomName}, c.Username)
				if err != nil {
					c.Send <- []byte("Failed create room " + roomName)
				} else {
					c.Send <- []byte("Room " + roomName + " created")
				}
			}
		} else if strings.HasPrefix(msg, "room:join:") {
			parts := strings.SplitN(msg, ":", 3)
			if len(parts) == 3 {
				roomName := parts[2]
				room, err := cc.roomUsecase.FindOne(&entity.Room{Name: roomName}, &pkg_model.Query{})
				if err != nil {
					c.Send <- []byte("Failed join room " + roomName)
				}
				user, err := cc.userUsecase.FindOne(&entity.User{Username: c.Username}, &pkg_model.Query{})
				if err != nil {
					c.Send <- []byte("Failed join room " + roomName)
				}
				err = cc.roomUserUsecase.Create(&entity.RoomUser{RoomID: room.ID, UserID: user.ID}, c.Username)
				if err != nil {
					c.Send <- []byte("Failed join room " + roomName)
				} else {
					c.Send <- []byte("Joined room " + roomName)
				}
			}
		} else if strings.HasPrefix(msg, "room:leave:") {
			parts := strings.SplitN(msg, ":", 3)
			if len(parts) == 3 {
				roomName := parts[2]
				room, err := cc.roomUsecase.FindOne(&entity.Room{Name: roomName}, &pkg_model.Query{})
				if err != nil {
					c.Send <- []byte("Failed leave room " + roomName)
				}
				user, err := cc.userUsecase.FindOne(&entity.User{Username: c.Username}, &pkg_model.Query{})
				if err != nil {
					c.Send <- []byte("Failed leave room " + roomName)
				}
				err = cc.roomUserUsecase.Delete(&entity.RoomUser{RoomID: room.ID, UserID: user.ID}, c.Username)
				if err != nil {
					c.Send <- []byte("Failed leave room " + roomName)
				} else {
					c.Send <- []byte("Left room " + roomName)
				}
			}
		} else if strings.HasPrefix(msg, "room:broadcast:") {
			parts := strings.SplitN(msg, ":", 4)
			if len(parts) == 4 {
				roomName := parts[2]
				room, err := cc.roomUsecase.FindOne(&entity.Room{Name: roomName}, &pkg_model.Query{Expand: []string{"RoomUsers.User"}})
				if err != nil {
					c.Send <- []byte("Failed broadcast to room " + roomName)
				}
				for _, ru := range room.RoomUsers {
					if ru.User.Username == c.Username {
						continue
					}
					if targetClient, ok := c.Hub.ClientsByUsername[ru.User.Username]; ok {
						targetClient.Send <- []byte("Broadcast from " + roomName + " (" + c.Username + "): " + parts[3])
						c.Send <- []byte("Successfully broadcast to room " + roomName + " (" + ru.User.Username + ")")
					}
				}
			}
		}
	})
}
