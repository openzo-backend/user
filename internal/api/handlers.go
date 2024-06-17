package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/service"
)

type Handler struct {
	userService     service.UserService
	userDataService service.UserDataService
}

func NewHandler(userService *service.UserService, userDataService *service.UserDataService) *Handler {
	return &Handler{userService: *userService, userDataService: *userDataService}
}

func (h *Handler) CreateUser(ctx *gin.Context) {
	var user service.CreateUserRequest
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUser, err := h.userService.CreateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

func (h *Handler) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := h.userService.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	user, err := h.userService.GetUserByEmail(ctx, email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *Handler) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUser(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}

func (h *Handler) UserSignIn(ctx *gin.Context) {
	var user service.UserSignInRequest
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.userService.UserSignIn(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetUserWithJWT(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")

	user, err := h.userService.GetUserWithJWT(ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// Add more handlers for other user operations (GetUser, UpdateUser, etc.)
func (h *Handler) CreateUserData(ctx *gin.Context) {
	var userData models.UserData
	if err := ctx.BindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdUserData, err := h.userDataService.CreateUserData(ctx, userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdUserData)
}

func (h *Handler) GetUserDataByID(ctx *gin.Context) {
	id := ctx.Param("id")

	userData, err := h.userDataService.GetUserDataByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userData)
}

func (h *Handler) UpdateUserData(ctx *gin.Context) {
	var userData models.UserData
	if err := ctx.BindJSON(&userData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUserData, err := h.userDataService.UpdateUserData(ctx, userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedUserData)
}

func (h *Handler) DeleteUserData(ctx *gin.Context) {
	id := ctx.Param("id")

	err := h.userDataService.DeleteUserData(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User data deleted successfully"})
}
