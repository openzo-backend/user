package service

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type OTPService interface {
	GenerateOTP(ctx *gin.Context, phoneNo string) (string, error)
	VerifyOTP(ctx *gin.Context, phone string, verificationId string, otp string, userId string) (string, error)
}

type otpService struct {
	otpRepository  repository.OTPRepository
	userRepository repository.UserRepository
}

func NewOTPService(otpRepository repository.OTPRepository,
	userRepository repository.UserRepository,
) OTPService {
	return &otpService{otpRepository: otpRepository, userRepository: userRepository}
}

func (s *otpService) GenerateOTP(ctx *gin.Context, phoneNo string) (string, error) {
	var otp models.OTP
	otp.Phone = phoneNo

	otp_number := 123456
	otp.HashedOTP = utils.HashNumberWithSecret(otp_number, "secret")

	generatedOTP, err := s.otpRepository.CreateOTP(otp)
	if err != nil {
		return "", err
	}

	return generatedOTP.ID, nil
}

func (s *otpService) VerifyOTP(ctx *gin.Context, phone string, verificationId string, otp string, userId string) (string, error) {
	_otp, err := s.otpRepository.GetOTPByID(verificationId)
	if err != nil {
		return "", err
	}

	if _otp.CreatedAt.Add(5 * time.Minute).Before(time.Now()) {
		return "", errors.New("OTP has expired")
	}

	otp_number, _ := strconv.Atoi(otp)

	if _otp.Phone != phone {
		return "", errors.New("invalid phone number")
	}

	if _otp.HashedOTP != utils.HashNumberWithSecret(otp_number, "secret") {
		return "", errors.New("invalid OTP")
	}

	// Delete the OTP from the database
	s.otpRepository.DeleteOTP(verificationId)

	user, err := s.userRepository.GetUserByMobile(phone)
	if err != nil || user.ID == "" {

		if userId != "" {
			user, err = s.userRepository.GetUserByID(userId)
			if err != nil {
				return "", err
			}
		} else {

			var newUser models.User
			newUser.Phone = phone
			newUser.CreatedAt = time.Now()
			createdUser, err := s.userRepository.CreateUser(newUser)
			if err != nil {
				return "", err
			}

			user = createdUser
		}

	} else if userId != "" && user.ID != userId {
		return "", errors.New("phone number already exists")
	}
	user.IsVerified = true
	user.Phone = phone

	_, err = s.userRepository.UpdateUser(user)
	if err != nil {
		return "", err
	}

	token, err := CreateJwtToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
