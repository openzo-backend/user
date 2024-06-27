package service

import (
	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/repository"
)

type UserService interface {

	//CRUD
	CreateUser(ctx *gin.Context, req models.User) (models.User, string, error)
	GetUserByID(ctx *gin.Context, id string) (models.User, error)
	GetUserByEmail(ctx *gin.Context, email string) (models.User, error)
	UpdateUser(ctx *gin.Context, req models.User) (models.User, error)

	//Authentication
	UserSignIn(ctx *gin.Context, req UserSignInRequest) (string, error)
	GetUserWithJWT(ctx *gin.Context, token string) (models.User, error)
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}

