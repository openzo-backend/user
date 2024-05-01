package service

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/utils"
)

type UserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *userService) UserSignIn(ctx *gin.Context, req UserSignInRequest) (string, error) {
	user, err := s.userRepository.GetUserByEmail(req.Email)

	if err != nil {
		return "", err
	}

	err = utils.CheckPasswordHash(req.Password, user.Password)
	if err != nil {
		return "", errors.New("invalid password")
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
