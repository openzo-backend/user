package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tanush-128/openzo_backend/user/internal/service"
)

type OTPHandler struct {
	otpService service.OTPService
}

func NewOTPHandler(otpService *service.OTPService) *OTPHandler {
	return &OTPHandler{otpService: *otpService}
}

type OTPRequest struct {
	PhoneNo string `json:"phone_no"`
}

type OTPVerifyRequest struct {
	PhoneNo        string `json:"phone_no"`
	OTP            string `json:"otp"`
	VerificationId string `json:"verification_id"`
}

func (h *OTPHandler) GenerateOTP(ctx *gin.Context) {
	var otpRequest OTPRequest
	if err := ctx.BindJSON(&otpRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verificationId, err := h.otpService.GenerateOTP(ctx, otpRequest.PhoneNo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"verification_id": verificationId})
}

func (h *OTPHandler) VerifyOTP(ctx *gin.Context) {
	var otpVerifyRequest OTPVerifyRequest
	if err := ctx.BindJSON(&otpVerifyRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	verified, err := h.otpService.VerifyOTP(ctx, otpVerifyRequest.PhoneNo, otpVerifyRequest.VerificationId, otpVerifyRequest.OTP)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, verified)
}
