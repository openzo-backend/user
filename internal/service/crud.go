package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
	// "github.com/tanush-128/openzo_backend/user/internal/repository"
)

type CreateUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (s *userService) CreateUser(ctx *gin.Context, req models.User) (models.User, string, error) {

	location, err := utils.GetLocation(*req.Latitude, *req.Longitude)
	if err != nil {
		return models.User{}, "", err
	}

	req.City = &location.Address.City
	req.State = &location.Address.State
	req.Country = &location.Address.Country
	req.Pincode = &location.Address.Postcode
	address := (location.Address.HouseNumber + ", " + location.Address.Road + ", " + location.Address.City + ", " + location.Address.State + ", " + location.Address.Country)
	req.Address = &address

	createdUser, err := s.userRepository.CreateUser(req)
	if err != nil {
		return models.User{}, "", err // Propagate error
	}

	token, err := CreateJwtToken(createdUser.ID)
	if err != nil {
		return models.User{}, "", err
	}

	return createdUser, token, nil
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
	location, err := utils.GetLocation(*req.Latitude, *req.Longitude)
	if err != nil {
		return models.User{}, err
	}

	req.City = &location.Address.City
	req.State = &location.Address.State
	req.Country = &location.Address.Country
	req.Pincode = &location.Address.Postcode
	address := (location.Address.HouseNumber + ", " + location.Address.Road + ", " + location.Address.City + ", " + location.Address.State + ", " + location.Address.Country)
	req.Address = &address
	updatedUser, err := s.userRepository.UpdateUser(req)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}
