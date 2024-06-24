package websocket

import (
	"net/http"

	"github.com/hanifmaliki/chat-app/internal/controller"
	"github.com/hanifmaliki/chat-app/internal/repository"
	"github.com/hanifmaliki/chat-app/internal/usecase"
	websocket "github.com/hanifmaliki/chat-app/pkg/websocket"

	"gorm.io/gorm"
)

func SetupRoutes(mux *http.ServeMux, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)
	userUsecase := usecase.NewUserUseCase(userRepo)
	userController := controller.NewUserController(userUsecase)

	roomRepo := repository.NewRoomRepository(db)
	roomUsecase := usecase.NewRoomUseCase(roomRepo)

	roomUserRepo := repository.NewRoomUserRepository(db)
	roomUserUsecase := usecase.NewRoomUserUseCase(roomUserRepo)

	hub := websocket.NewHub()
	go hub.Run()
	chatController := controller.NewChatController(userUsecase, roomUsecase, roomUserUsecase, hub)

	mux.HandleFunc("/health", controller.HealthCheckHandler)
	mux.HandleFunc("/register", userController.Register)
	mux.HandleFunc("/login", userController.Login)
	mux.HandleFunc("/chat/ws", chatController.HandleWebSocket)
}
