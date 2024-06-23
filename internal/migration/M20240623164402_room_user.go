package migration

import (
	"github.com/hanifmaliki/chat-app/pkg/model"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var M20240623164402_room_user = gormigrate.Migration{
	ID: "M20240623164402_room_user",
	Migrate: func(tx *gorm.DB) error {
		type RoomUser struct {
			model.Base

			RoomID uint `gorm:"index:uidx_room_users,unique,priority:10"`
			UserID uint `gorm:"index:uidx_room_users,unique,priority:9"`
		}

		return tx.Migrator().CreateTable(&RoomUser{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("room_users")
	},
}
