package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	// "github.com/tanush-128/openzo_backend/user/internal/repository"
)

type CreateUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (s *userService) CreateUser(ctx *gin.Context, req CreateUserRequest) (models.User, error) {


	user := models.User{
		Phone: req.Phone,
	}

	createdUser, err := s.userRepository.CreateUser(user)
	if err != nil {
		return models.User{}, err // Propagate error
	}

	return createdUser, nil
}

func (s *userService) GetUserByID(ctx *gin.Context, id string) (models.User, error) {
	user, err := s.userRepository.GetUserByID(id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *userService) GetUserByEmail(ctx *gin.Context, email string) (models.User, error) {
	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *userService) UpdateUser(ctx *gin.Context, req models.User) (models.User, error) {

	updatedUser, err := s.userRepository.UpdateUser(req)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}
