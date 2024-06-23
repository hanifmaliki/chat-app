package entity

import "github.com/hanifmaliki/chat-app/pkg/model"

type RoomUser struct {
	model.Base

	RoomID uint `json:"room_id"`
	UserID uint `json:"user_id"`

	Room *Room `json:"room"`
	User *User `json:"user"`
}
