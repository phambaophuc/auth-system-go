package middleware

import (
	"auth-system/internal/infrastructure/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *security.JWTManager
}

func NewAuthMiddleware(jwtManager *security.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		// Check if starts with "Bearer "
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		if claims.Type != "access" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid token type",
			})
			ctx.Abort()
			return
		}

		// Store user information in context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)
		ctx.Next()
	}
}

func (m *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.Next()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.Next()
			return
		}

		tokenString := parts[1]
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		if claims.Type == "access" {
			ctx.Set("user_id", claims.UserID)
			ctx.Set("user_email", claims.Email)
		}

		ctx.Next()
	}
}
