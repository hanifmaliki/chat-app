package usecase

// import (
// 	"github.com/hanifmaliki/chat-app/internal/entity"
// 	"github.com/hanifmaliki/chat-app/internal/repository"
// 	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"
// )

// type ChatUseCase struct {
// 	userUsecase     *UserUseCase
// 	roomUsecase     *RoomUseCase
// 	roomUserUsecase *RoomUserUseCase
// }

// func NewChatUseCase(userUsecase *UserUseCase, roomUsecase *RoomUseCase, roomUserUsecase *RoomUserUseCase) *ChatUseCase {
// 	return &ChatUseCase{
// 		userUsecase:     userUsecase,
// 		roomUsecase:     roomUsecase,
// 		roomUserUsecase: roomUserUsecase,
// 	}
// }

// func (uc *ChatUseCase) Create(data *entity.Room, by string) error {
// 	return uc.repo.Create(data, by)
// }
