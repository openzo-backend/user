package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"gorm.io/gorm"
)

type UserDataRepository interface {
	CreateUserData(user models.UserData) (models.UserData, error)
	GetUserDataByID(id string) (models.UserData, error)
	UpdateUserData(user models.UserData) (models.UserData, error)
	DeleteUserData(id string) error
	// Add more methods for other user operations (GetUserByEmail, UpdateUser, etc.)

}

type userDataRepository struct {
	db *gorm.DB
}

func NewUserDataRepository(db *gorm.DB) UserDataRepository {

	return &userDataRepository{db: db}
}

func (r *userDataRepository) CreateUserData(user models.UserData) (models.UserData, error) {
	user.Id = uuid.New().String()
	tx := r.db.Create(&user)
	if tx.Error != nil {
		return models.UserData{}, tx.Error
	}

	return user, nil
}

func (r *userDataRepository) GetUserDataByID(id string) (models.UserData, error) {
	var user models.UserData
	tx := r.db.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return models.UserData{}, tx.Error
	}

	return user, nil
}

func (r *userDataRepository) UpdateUserData(user models.UserData) (models.UserData, error) {

	tx := r.db.Save(&user)
	if tx.Error != nil {
		return models.UserData{}, tx.Error
	}

	return user, nil
}

func (r *userDataRepository) DeleteUserData(id string) error {
	tx := r.db.Where("id = ?", id).Delete(&models.UserData{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
