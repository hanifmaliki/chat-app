package migration

import (
	"github.com/hanifmaliki/chat-app/pkg/model"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var M20240623164400_user = gormigrate.Migration{
	ID: "M20240623164400_user",
	Migrate: func(tx *gorm.DB) error {
		type User struct {
			model.Base

			Username string `gorm:"index:,unique"`
			Password string
		}

		return tx.Migrator().CreateTable(&User{})
	},
	Rollback: func(tx *gorm.DB) error {
		return tx.Migrator().DropTable("users")
	},
}
