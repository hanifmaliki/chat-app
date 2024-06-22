package usecase

import (
	"github.com/hanifmaliki/chat-app/internal/model"
	"github.com/hanifmaliki/chat-app/internal/repository"
)

type ChatService struct {
	messageRepo repository.MessageRepository
}

func NewChatService(repo repository.MessageRepository) *ChatService {
	return &ChatService{messageRepo: repo}
}

func (s *ChatService) SaveMessage(room, user, content string) error {
	msg := &model.Message{Room: room, User: user, Content: content}
	return s.messageRepo.Save(msg)
}

func (s *ChatService) GetMessages(room string) ([]model.Message, error) {
	return s.messageRepo.GetMessages(room)
}
