package repository

import (
	"github.com/hanifmaliki/chat-app/internal/model"

	"gorm.io/gorm"
)

type MessageRepository interface {
	Save(message *model.Message) error
	GetMessages(room string) ([]model.Message, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewGormMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Save(message *model.Message) error {
	return r.db.Create(message).Error
}

func (r *messageRepository) GetMessages(room string) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.Where("room = ?", room).Order("timestamp").Find(&messages).Error
	return messages, err
}
