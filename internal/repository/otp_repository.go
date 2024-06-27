package repository

import (
	"github.com/google/uuid"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"gorm.io/gorm"
)

// type OTP struct {
// 	ID        string `gorm:"primaryKey"`
// 	HashedOTP string
// 	CreatedAt time.Time `gorm:"autoCreateTime"`
// }

type OTPRepository interface {
	CreateOTP(otp models.OTP) (models.OTP, error)
	GetOTPByID(id string) (models.OTP, error)
	DeleteOTP(id string) error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {

	return &otpRepository{db: db}
}

func (r *otpRepository) CreateOTP(otp models.OTP) (models.OTP, error) {
	otp.ID = uuid.New().String()
	tx := r.db.Create(&otp)
	if tx.Error != nil {
		return models.OTP{}, tx.Error
	}

	return otp, nil
}

func (r *otpRepository) GetOTPByID(id string) (models.OTP, error) {
	var otp models.OTP
	tx := r.db.Where("id = ?", id).First(&otp)
	if tx.Error != nil {
		return models.OTP{}, tx.Error
	}

	return otp, nil
}


func (r *otpRepository) DeleteOTP(id string) error {
	tx := r.db.Where("id = ?", id).Delete(&models.OTP{})
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
