package usecase

import (
	"github.com/hanifmaliki/chat-app/internal/entity"
	"github.com/hanifmaliki/chat-app/internal/repository"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"
)

type RoomUseCase struct {
	repo repository.RoomRepository
}

func NewRoomUseCase(repo repository.RoomRepository) *RoomUseCase {
	return &RoomUseCase{repo: repo}
}

func (uc *RoomUseCase) Create(data *entity.Room, by string) error {
	return uc.repo.Create(data, by)
}

func (uc *RoomUseCase) Delete(conds *entity.Room, by string) error {
	return uc.repo.Delete(conds, by)
}

func (uc *RoomUseCase) Find(conds *entity.Room, query *pkg_model.Query) ([]*entity.Room, error) {
	return uc.repo.Find(conds, query)
}

func (uc *RoomUseCase) FindOne(conds *entity.Room, query *pkg_model.Query) (*entity.Room, error) {
	return uc.repo.FindOne(conds, query)
}
