package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tanush-128/openzo_backend/user/internal/models"
)

type UserSignInRequest struct {
	Mobile string `json:"mobile"`
	Otp    string `json:"otp"`
}

func (s *userService) UserSignIn(ctx *gin.Context, req UserSignInRequest) (string, error) {
	// Validate user data (implement validation logic)
	if req.Mobile == "" {
		return "", errors.New("invalid request")
	}

	// Get user by mobile
	user, err := s.userRepository.GetUserByMobile(req.Mobile)
	if err != nil {
		return "", err
	}

	// Create JWT token
	token, err := CreateJwtToken(user.ID)
	if err != nil {
		return token, err
	}

	return token, nil
}

// GetUserWithJWT(ctx *gin.Context, token string) (models.User, error)
func (s *userService) GetUserWithJWT(ctx *gin.Context, token string) (models.User, error) {
	_user := ctx.MustGet("user").(map[string]interface{})
	_id := _user["user_id"].(string)

	user, err := s.userRepository.GetUserByID(_id)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func CreateJwtToken(id string) (string, error) {
	claims := jwt.MapClaims{}

	claims["authorized"] = true
	claims["user_id"] = id
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
