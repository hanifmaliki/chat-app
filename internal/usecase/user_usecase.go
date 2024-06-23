package usecase

import (
	"errors"

	"github.com/hanifmaliki/chat-app/internal/entity"
	"github.com/hanifmaliki/chat-app/internal/model"
	"github.com/hanifmaliki/chat-app/internal/repository"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) Create(data *entity.User, by string) error {
	return uc.repo.Create(data, by)
}

func (uc *UserUseCase) Delete(conds *entity.User, by string) error {
	return uc.repo.Delete(conds, by)
}

func (uc *UserUseCase) Find(conds *entity.User, query *pkg_model.Query) ([]*entity.User, error) {
	return uc.repo.Find(conds, query)
}

func (uc *UserUseCase) FindOne(conds *entity.User, query *pkg_model.Query) (*entity.User, error) {
	return uc.repo.FindOne(conds, query)
}

func (uc *UserUseCase) Register(data *model.Credential) (*entity.User, error) {
	user := entity.User{Username: data.Username, Password: data.Password}
	err := uc.repo.Create(&user, "")
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (uc *UserUseCase) Login(data *model.Credential) (*entity.User, error) {
	user, err := uc.repo.FindOne(&entity.User{Username: data.Username}, &pkg_model.Query{})
	if err != nil {
		return nil, err
	}

	if user.Password != data.Password {
		return nil, errors.New("invalid password")
	}

	return user, nil
}
