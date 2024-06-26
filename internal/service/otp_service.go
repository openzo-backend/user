package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid_v4"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
)

type OTPService interface {
	GenerateOTP(ctx *gin.Context, phoneNo string) (string, error)
	VerifyOTP(ctx *gin.Context, phone string, verificationId string, otp string) (bool, error)
}

type otpService struct {
	otpRepository repository.OTPRepository
}

func NewOTPService(otpRepository repository.OTPRepository) OTPService {
	return &otpService{otpRepository: otpRepository}
}

func hashNumberWithSecret(number int, secret string) string {
	// Convert the number to a byte slice
	numberBytes := []byte(fmt.Sprintf("%d", number))

	// Create a new HMAC by defining the hash type and the key (secret)
	hmac := hmac.New(sha256.New, []byte(secret))

	// Write the number bytes to the HMAC object
	hmac.Write(numberBytes)

	// Compute the HMAC hash
	hash := hmac.Sum(nil)

	// Encode the hash to a hexadecimal string and return
	return hex.EncodeToString(hash)
}

func (s *otpService) GenerateOTP(ctx *gin.Context, phoneNo string) (string, error) {
	var otp models.OTP
	otp.Phone = phoneNo

	otp_number := 123456
	otp.HashedOTP = hashNumberWithSecret(otp_number, "secret")

	generatedOTP, err := s.otpRepository.CreateOTP(otp)
	if err != nil {
		return "", err
	}

	return generatedOTP.ID, nil
}

func (s *otpService) VerifyOTP(ctx *gin.Context, phone string, verificationId string, otp string) (bool, error) {
	_otp, err := s.otpRepository.GetOTPByID(verificationId)
	if err != nil {
		return false, err
	}

	if _otp.CreatedAt.Add(5 * time.Minute).Before(time.Now()) {
		return false, errors.New("OTP has expired")
	}

	otp_number, _ := strconv.Atoi(otp)

	if _otp.Phone != phone {
		return false, errors.New("invalid phone number")
	}

	if _otp.HashedOTP != hashNumberWithSecret(otp_number, "secret") {
		return true, errors.New("invalid OTP")
	}

	return true, nil
}
