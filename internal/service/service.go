package service

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type UserService interface {

	//CRUD
	CreateUser(ctx *gin.Context, req models.User) (models.User, string, error)
	GetUserByID(ctx *gin.Context, id string) (models.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (models.User, error)
	UpdateUser(ctx *gin.Context, req models.User) (models.User, error)

	//Authentication
	UserSignIn(ctx *gin.Context, req UserSignInRequest) (string, error)
	GetUserWithJWT(ctx *gin.Context, token string) (models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

type CreateUserRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (s *userService) CreateUser(ctx *gin.Context, req models.User) (models.User, string, error) {

	if req.Latitude != nil && req.Longitude != nil {
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
	}

	hashedPassword, err := utils.HashPassword(*req.Password)
	if err != nil {
		return models.User{}, "", err
	}

	req.Password = &hashedPassword

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

	var user models.User

	user, err := s.userRepository.GetUserByID(req.ID)
	if err != nil {
		return models.User{}, err
	}

	if (req.Latitude != nil && req.Longitude != nil) && (user.Pincode == nil) {

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
	} else if user.Pincode != req.Pincode {
		locations, err := utils.GetLocationByPincode(*req.Pincode)
		if err != nil {
			return models.User{}, err
		}
		location := locations[0]

		req.Address = &location.DisplayName
		req.Latitude = &location.Lat
		req.Longitude = &location.Lon
		city := strings.Split(location.DisplayName, ", ")[0]
		req.City = &city
		req.State = &strings.Split(location.DisplayName, ", ")[3]

	}
	req.Password = user.Password
	req.Role = user.Role
	updatedUser, err := s.userRepository.UpdateUser(req)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}
