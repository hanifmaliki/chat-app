package migration

import (
	"github.com/hanifmaliki/chat-app/pkg/model"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var M20240623164401_room = gormigrate.Migration{
	ID: "M20240623164401_room",
	Migrate: func(tx *gorm.DB) error {
		type Room struct {
			model.Base

			Name string `gorm:"index:,unique"`
		}

		return tx.Migrator().CreateTable(&Room{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("rooms")
	},
}
