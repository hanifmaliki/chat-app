package usecase

import (
	"github.com/hanifmaliki/chat-app/internal/model"
	"github.com/hanifmaliki/chat-app/internal/repository"
)

type ChatUseCase struct {
	messageRepo repository.MessageRepository
}

func NewChatService(repo repository.MessageRepository) *ChatUseCase {
	return &ChatUseCase{messageRepo: repo}
}

func (s *ChatUseCase) SaveMessage(room, user, content string) error {
	msg := &model.Message{Room: room, User: user, Content: content}
	return s.messageRepo.Save(msg)
}

func (s *ChatUseCase) GetMessages(room string) ([]model.Message, error) {
	return s.messageRepo.GetMessages(room)
}
