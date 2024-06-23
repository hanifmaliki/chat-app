package usecase

import (
	"github.com/hanifmaliki/chat-app/internal/entity"
	"github.com/hanifmaliki/chat-app/internal/repository"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"
)

type RoomUserUseCase struct {
	repo repository.RoomUserRepository
}

func NewRoomUserUseCase(repo repository.RoomUserRepository) *RoomUserUseCase {
	return &RoomUserUseCase{repo: repo}
}

func (uc *RoomUserUseCase) Create(data *entity.RoomUser, by string) error {
	return uc.repo.Create(data, by)
}

func (uc *RoomUserUseCase) Delete(conds *entity.RoomUser, by string) error {
	return uc.repo.Delete(conds, by)
}

func (uc *RoomUserUseCase) Find(conds *entity.RoomUser, query *pkg_model.Query) ([]*entity.RoomUser, error) {
	return uc.repo.Find(conds, query)
}

func (uc *RoomUserUseCase) FindOne(conds *entity.RoomUser, query *pkg_model.Query) (*entity.RoomUser, error) {
	return uc.repo.FindOne(conds, query)
}
