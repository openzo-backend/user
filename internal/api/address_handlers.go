package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/models"
	"github.com/tanush-128/openzo_backend/user/internal/service"
)

type AddressHandler struct {
	addressService service.AddressService
}

func NewAddressHandler(addressService *service.AddressService) *AddressHandler {
	return &AddressHandler{addressService: *addressService}
}

func (h *AddressHandler) CreateAddress(ctx *gin.Context) {
	var address models.Address
	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdAddress, err := h.addressService.CreateAddress(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdAddress)
}

func (h *AddressHandler) GetAddressByID(ctx *gin.Context) {
	id := ctx.Param("id")

	address, err := h.addressService.GetAddressByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (h *AddressHandler) GetAddressesByUserID(ctx *gin.Context) {
	user_id := ctx.Param("user_id")

	address, err := h.addressService.GetAddressesByUserId(ctx, user_id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, address)
}

func (h *AddressHandler) UpdateAddress(ctx *gin.Context) {
	var address models.Address
	if err := ctx.BindJSON(&address); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedAddress, err := h.addressService.UpdateAddress(ctx, address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedAddress)
}
