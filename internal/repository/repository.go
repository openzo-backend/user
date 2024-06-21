package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user models.User) (models.User, error)
	GetUserByID(id string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByMobile(mobile string) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	// Add more methods for other user operations (GetUserByEmail, UpdateUser, etc.)

}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {

	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user models.User) (models.User, error) {
	var user2 models.User
	r.db.Find(
		&user2, "phone = ?", user.Phone,
	)
	if user2.Phone != "" {
		return models.User{}, errors.New("user already exists")
	}

	user.ID = uuid.New().String()

	tx := r.db.Create(&user)

	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}

func (r *userRepository) GetUserByID(id string) (models.User, error) {
	var user models.User
	tx := r.db.Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	tx := r.db.Where("email = ?", email).First(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}

func (r *userRepository) GetUserByMobile(mobile string) (models.User, error) {
	var user models.User
	tx := r.db.Where("phone = ?", mobile).First(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}

func (r *userRepository) UpdateUser(user models.User) (models.User, error) {
	tx := r.db.Save(&user)
	if tx.Error != nil {
		return models.User{}, tx.Error
	}

	return user, nil
}

// Implement other repository methods (GetUserByID, GetUserByEmail, UpdateUser, etc.) with proper error handling
