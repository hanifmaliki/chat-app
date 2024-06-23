package db

import (
	"github.com/hanifmaliki/chat-app/internal/migration"
	"github.com/hanifmaliki/chat-app/pkg/util"

	gormigrate "github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewGormDB() (*gorm.DB, error) {
	dbURL := util.GetEnv("DATABASE_URL", "chat.db")
	db, err := gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	m := gormigrate.New(db, gormigrate.DefaultOptions, migration.Migrations)

	if err := m.Migrate(); err != nil {
		return nil, err
	}

	return db, nil
}
