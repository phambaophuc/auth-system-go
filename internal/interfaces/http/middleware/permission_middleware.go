package middleware

import (
	"auth-system/internal/domain/services"
	"auth-system/pkg/errors"
	"auth-system/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

type PermissionMiddleware struct {
	permissionService services.PermissionService
}

func NewPermissionMiddleware(permissionService services.PermissionService) *PermissionMiddleware {
	return &PermissionMiddleware{
		permissionService: permissionService,
	}
}

func (m *PermissionMiddleware) RequirePermission(resource, action string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")
		if !exists {
			utils.ErrorResponse(ctx, errors.NewUnauthorizedError("Authentication required"))
			ctx.Abort()
			return
		}

		userIDUint, ok := userID.(uint)
		if !ok {
			utils.ErrorResponse(ctx, errors.NewInternalServerError("Invalid user id"))
			ctx.Abort()
			return
		}

		hasPermission, err := m.permissionService.CheckPermission(userIDUint, resource, action)
		if err != nil {
			utils.ErrorResponse(ctx, errors.NewInternalServerError(fmt.Sprintf("failed to check permission: %v", err)))
			ctx.Abort()
			return
		}

		if !hasPermission {
			utils.ErrorResponse(ctx, errors.NewForbiddenError("Insufficient permissions"))
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func (m *PermissionMiddleware) RequireAnyPermission(permissions []struct{ Resource, Action string }) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := ctx.Get("user_id")
		if !exists {
			utils.ErrorResponse(ctx, errors.NewUnauthorizedError("Authentication required"))
			ctx.Abort()
			return
		}

		userIDUint, ok := userID.(uint)
		if !ok {
			utils.ErrorResponse(ctx, errors.NewInternalServerError("Invalid user id"))
			ctx.Abort()
			return
		}

		for _, perm := range permissions {
			hasPermission, err := m.permissionService.CheckPermission(userIDUint, perm.Resource, perm.Action)
			if err != nil {
				utils.ErrorResponse(ctx, errors.NewInternalServerError(fmt.Sprintf("Failed to check permission: %v", err)))
				ctx.Abort()
				return
			}

			if hasPermission {
				ctx.Next()
				return
			}
		}

		utils.ErrorResponse(ctx, errors.NewForbiddenError("Insufficient permissions"))
		ctx.Abort()
	}
}
