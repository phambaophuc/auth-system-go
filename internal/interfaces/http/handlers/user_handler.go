package handlers

import (
	"auth-system/internal/application/dto"
	"auth-system/internal/domain/services"
	"auth-system/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
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

	profile, err := h.userService.GetProfile(userIDUint)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile retrieved successfully", profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
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

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	profile, err := h.userService.UpdateProfile(userIDUint, &req)
	if err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Profile updated successfully", profile)
}

func (h *UserHandler) ChangePassword(c *gin.Context) {
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

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	if err := h.userService.ChangePassword(userIDUint, &req); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Password changed successfully", nil)
}

func (h *UserHandler) AssignRole(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request data")
		return
	}

	if err := h.userService.AssignRole(uint(userID), req.RoleID); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Role assigned successfully", nil)
}

func (h *UserHandler) RemoveRole(c *gin.Context) {
	userIDParam := c.Param("id")
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid user ID")
		return
	}

	roleIDParam := c.Param("roleId")
	roleID, err := strconv.ParseUint(roleIDParam, 10, 32)
	if err != nil {
		utils.ValidationErrorResponse(c, "Invalid role ID")
		return
	}

	if err := h.userService.RemoveRole(uint(userID), uint(roleID)); err != nil {
		utils.ErrorResponse(c, err)
		return
	}

	utils.SuccessResponse(c, "Role removed successfully", nil)
}
