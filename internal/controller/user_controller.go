package controller

import (
	"encoding/json"
	"net/http"

	"github.com/hanifmaliki/chat-app/internal/entity"
	"github.com/hanifmaliki/chat-app/internal/model"
	"github.com/hanifmaliki/chat-app/internal/usecase"
	pkg_model "github.com/hanifmaliki/chat-app/pkg/model"

	"github.com/rs/zerolog/log"
)

type UserController struct {
	userUsecase *usecase.UserUseCase
}

func NewUserController(userUsecase *usecase.UserUseCase) *UserController {
	return &UserController{
		userUsecase: userUsecase,
	}
}

func (cc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req model.Credential

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		response := pkg_model.Response[any]{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request payload",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := cc.userUsecase.Register(&req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register user")
		response := pkg_model.Response[any]{
			Code:    http.StatusInternalServerError,
			Success: false,
			Message: "Failed to register user",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := pkg_model.Response[*entity.User]{
		Code:    http.StatusCreated,
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	}
	json.NewEncoder(w).Encode(response)
}

func (cc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req model.Credential

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		response := pkg_model.Response[any]{
			Code:    http.StatusBadRequest,
			Success: false,
			Message: "Invalid request payload",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	user, err := cc.userUsecase.Login(&req)
	if err != nil {
		log.Error().Err(err).Msg("Invalid username or password")
		response := pkg_model.Response[any]{
			Code:    http.StatusUnauthorized,
			Success: false,
			Message: "Invalid username or password",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := pkg_model.Response[*entity.User]{
		Code:    http.StatusOK,
		Success: true,
		Message: "Login successful",
		Data:    user,
	}
	json.NewEncoder(w).Encode(response)
}
