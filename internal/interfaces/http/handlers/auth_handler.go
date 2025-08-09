package handlers

import (
	"auth-system/internal/application/dto"
	"auth-system/internal/domain/services"
	"auth-system/pkg/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	response, err := h.authService.Login(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Login successful", response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	response, err := h.authService.Register(&req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.CreatedResponse(c, "Registration successful", response)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Token refreshed successfully", response)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ValidationErrorResponse(c, "User not found in context")
		return
	}

	userIDUint, ok := userID.(uint)
	if !ok {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	if err := h.authService.Logout(userIDUint); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Logout successful", nil)
}
