package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type AddressService interface {

	//CRUD
	CreateAddress(ctx *gin.Context, req models.Address) (models.Address, error)
	GetAddressByID(ctx *gin.Context, id string) (models.Address, error)
	GetAddressesByUserId(ctx *gin.Context, user_id string) ([]models.Address, error)
	UpdateAddress(ctx *gin.Context, req models.Address) (models.Address, error)
}
type addressService struct {
	addressRepository repository.AddressRepository
}

func NewAddressService(addressRepository repository.AddressRepository) AddressService {
	return &addressService{addressRepository: addressRepository}
}

func (s *addressService) CreateAddress(ctx *gin.Context, req models.Address) (models.Address, error) {

	if req.Latitude != "" && req.Longitude != "" {
		location, err := utils.GetLocation(req.Latitude, req.Longitude)
		if err != nil {
			return models.Address{}, err
		}

		req.City = location.Address.City
		req.State = location.Address.State

		req.Pincode = location.Address.Postcode
		address := (location.Address.HouseNumber + ", " + location.Address.Road + ", " + location.Address.City + ", " + location.Address.State + ", " + location.Address.Country)
		req.Address = address
	}

	createdAddress, err := s.addressRepository.CreateAddress(req)
	if err != nil {
		return models.Address{}, err // Propagate error
	}

	return createdAddress, nil
}

func (s *addressService) GetAddressByID(ctx *gin.Context, id string) (models.Address, error) {
	address, err := s.addressRepository.GetAddressByID(id)
	if err != nil {
		return models.Address{}, err
	}

	return address, nil
}

func (s *addressService) GetAddressesByUserId(ctx *gin.Context, user_id string) ([]models.Address, error) {
	addresses, err := s.addressRepository.GetAddressesByUserID(user_id)
	if err != nil {
		return nil, err
	}

	return addresses, nil
}

func (s *addressService) UpdateAddress(ctx *gin.Context, req models.Address) (models.Address, error) {

	updatedAddress, err := s.addressRepository.UpdateAddress(req)
	if err != nil {
		return models.Address{}, err
	}

	return updatedAddress, nil
}
