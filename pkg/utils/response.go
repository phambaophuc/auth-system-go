package utils

import (
	"auth-system/internal/application/dto"
	"auth-system/pkg/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := err.(*errors.AppError); ok {
		c.JSON(appErr.Code, dto.APIResponse{
			Success: false,
			Message: appErr.Message,
			Error:   appErr.Message,
		})
		return
	}

	// Default to internal server error
	c.JSON(http.StatusInternalServerError, dto.APIResponse{
		Success: false,
		Message: "Internal server error",
		Error:   err.Error(),
	})
}

func ValidationErrorResponse(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, dto.APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	})
}
