package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"gorm.io/gorm"
)

type AddressRepository interface {
	CreateAddress(address models.Address) (models.Address, error)
	GetAddressByID(id string) (models.Address, error)
	GetAddressesByUserID(user_id string) ([]models.Address, error)
	UpdateAddress(address models.Address) (models.Address, error)
}

type addressRepository struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) AddressRepository {

	return &addressRepository{db: db}
}

func (r *addressRepository) CreateAddress(address models.Address) (models.Address, error) {

	address.ID = uuid.New().String()

	tx := r.db.Create(&address)

	if tx.Error != nil {
		return models.Address{}, tx.Error
	}

	return address, nil
}

func (r *addressRepository) GetAddressByID(id string) (models.Address, error) {
	var address models.Address
	tx := r.db.Where("id = ?", id).First(&address)
	if tx.Error != nil {
		return models.Address{}, tx.Error
	}

	return address, nil
}

func (r *addressRepository) UpdateAddress(address models.Address) (models.Address, error) {
	tx := r.db.Save(&address)
	if tx.Error != nil {
		return models.Address{}, tx.Error
	}

	return address, nil
}

func (r *addressRepository) GetAddressesByUserID(user_id string) ([]models.Address, error) {
	var addresses []models.Address
	tx := r.db.Where("user_id = ?", user_id).Find(&addresses)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return addresses, nil
}
