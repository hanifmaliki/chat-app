package entity

import "github.com/hanifmaliki/chat-app/pkg/model"

type User struct {
	model.Base

	Username string `json:"username"`
	Password string `json:"password"`

	RoomUsers []*RoomUser `json:"room_users"`
}
