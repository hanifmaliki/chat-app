package entity

import "github.com/hanifmaliki/chat-app/pkg/model"

type Room struct {
	model.Base

	Name string `json:"name"`

	RoomUsers []*RoomUser `json:"room_users"`
}
