package routes

import (
	"auth-system/internal/interfaces/http/handlers"
	"auth-system/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler    *handlers.AuthHandler
	userHandler    *handlers.UserHandler
	authMiddleware *middleware.AuthMiddleware
	permMiddleware *middleware.PermissionMiddleware
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	authMiddleware *middleware.AuthMiddleware,
	permMiddleware *middleware.PermissionMiddleware,
) *Router {
	return &Router{
		authHandler:    authHandler,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
		permMiddleware: permMiddleware,
	}
}

func (r *Router) SetupRoutes() *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"status":  "Healthy",
			"message": "Auth service is running",
		})
	})

	api := router.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	{
		auth.POST("/login", r.authHandler.Login)
		auth.POST("/register", r.authHandler.Register)
		auth.POST("/refresh", r.authHandler.RefreshToken)
		auth.POST("/logout", r.authMiddleware.RequireAuth(), r.authHandler.Logout)
	}

	// Protected routes
	user := api.Group("/user")
	user.Use(r.authMiddleware.RequireAuth())
	{
		user.GET("/profile", r.userHandler.GetProfile)
		user.PUT("/profile", r.userHandler.UpdateProfile)
		user.POST("/change-password", r.userHandler.ChangePassword)
	}

	// Admin routes
	admin := api.Group("/admin")
	admin.Use(r.authMiddleware.RequireAuth())
	admin.Use(r.permMiddleware.RequirePermission("users", "write"))
	{
		admin.POST("/users/:id/roles", r.userHandler.AssignRole)
		admin.DELETE("/users/:id/roles/:roleId", r.userHandler.RemoveRole)
	}

	return router
}
