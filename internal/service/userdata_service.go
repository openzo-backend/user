package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type UserDataService interface {

	//CRUD
	CreateUserData(ctx *gin.Context, req models.UserData) (models.UserData, error)
	GetUserDataByID(ctx *gin.Context, id string) (models.UserData, error)
	UpdateUserData(ctx *gin.Context, req models.UserData) (models.UserData, error)
	DeleteUserData(ctx *gin.Context, id string) error
	GetUserDataByUserID(ctx *gin.Context, id string) (models.UserData, error)
}

type userDataService struct {
	userDataRepository repository.UserDataRepository
}

func NewUserDataService(userDataRepository repository.UserDataRepository) UserDataService {
	return &userDataService{userDataRepository: userDataRepository}
}

func (s *userDataService) CreateUserData(ctx *gin.Context, req models.UserData) (models.UserData, error) {

	location, err := utils.GetLocation(req.Latitude, req.Longitude)
	if err != nil {
		return models.UserData{}, err
	}

	req.City = location.Address.City
	req.State = location.Address.State
	req.Country = location.Address.Country
	req.Pincode = location.Address.Postcode
	req.Address = location.Address.HouseNumber + ", " + location.Address.Road + ", " + location.Address.City + ", " + location.Address.State + ", " + location.Address.Country

	return s.userDataRepository.CreateUserData(req)
}

func (s *userDataService) GetUserDataByID(ctx *gin.Context, id string) (models.UserData, error) {

	return s.userDataRepository.GetUserDataByID(id)
}

func (s *userDataService) GetUserDataByUserID(ctx *gin.Context, id string) (models.UserData, error) {
	return s.userDataRepository.GetUserDataByUserID(id)
}

func (s *userDataService) UpdateUserData(ctx *gin.Context, req models.UserData) (models.UserData, error) {

	return s.userDataRepository.UpdateUserData(req)
}

func (s *userDataService) DeleteUserData(ctx *gin.Context, id string) error {

	return s.userDataRepository.DeleteUserData(id)
}
