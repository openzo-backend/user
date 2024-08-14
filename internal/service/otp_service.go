package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/config"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type OTPService interface {
	GenerateOTP(ctx *gin.Context, phoneNo string) (string, error)
	VerifyOTP(ctx *gin.Context, phone string, verificationId string, otp string, userId string) (string, error)
	SendOTP(phoneNo string, otp string)
}

type otpService struct {
	otpRepository  repository.OTPRepository
	userRepository repository.UserRepository
	cfg            *config.Config
}

func NewOTPService(otpRepository repository.OTPRepository,
	userRepository repository.UserRepository,
	cfg *config.Config,
) OTPService {
	return &otpService{otpRepository: otpRepository, userRepository: userRepository, cfg: cfg}
}

func (s *otpService) GenerateOTP(ctx *gin.Context, phoneNo string) (string, error) {
	var otp models.OTP
	otp.Phone = phoneNo

	otp_number := generatedRandomOTP()
	fmt.Println(otp_number)
	go s.SendOTP(phoneNo, strconv.Itoa(otp_number))
	otp.HashedOTP = utils.HashNumberWithSecret(otp_number, "secret")

	generatedOTP, err := s.otpRepository.CreateOTP(otp)
	if err != nil {
		return "", err
	}

	return generatedOTP.ID, nil
}

func generatedRandomOTP() int {
	random := 1000 + rand.Intn(9000)

	return random
	// return strconv.Itoa(random)
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

	}
	// else if userId != "" && user.ID != userId {
	// 	return "", errors.New("phone number already exists")
	// }
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

func (s *otpService) SendOTP(phoneNo string, otp string) {
	// url := "https://2factor.in/API/V1/fae85dd6-50a7-11ef-8b60-0200cd936042/SMS/+919999999999/12345/OTP1"
	url := "https://2factor.in/API/V1/" + s.cfg.SMS_API_KEY + "/SMS/" + phoneNo + "/" + otp + "/OTP 1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
