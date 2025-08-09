package main

import (
	"auth-system/internal/application/services"
	"auth-system/internal/config"
	"auth-system/internal/infrastructure/database"
	"auth-system/internal/infrastructure/repositories"
	"auth-system/internal/infrastructure/security"
	"auth-system/internal/interfaces/http/handlers"
	"auth-system/internal/interfaces/http/middleware"
	"auth-system/internal/interfaces/http/routes"
	"log"
	"net/http"
)

func main() {
	cfg := config.Load()

	db, err := database.NewConnection(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := database.AutoMigrate(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	jwtManager, err := security.NewJWTManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenTTL,
		cfg.JWT.RefreshTokenTTL,
	)
	if err != nil {
		log.Fatal("Failed to initialize JWT manager:", err)
	}

	passwordManager := security.NewPasswordManager()

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	roleRepo := repositories.NewRoleRepository(db)
	permissionRepo := repositories.NewPermissionRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, roleRepo, jwtManager, passwordManager)
	userService := services.NewUserService(userRepo, roleRepo, passwordManager)
	permissionService := services.NewPermissionService(permissionRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtManager)
	permMiddleware := middleware.NewPermissionMiddleware(permissionService)

	// Setup routes
	router := routes.NewRouter(
		authHandler,
		userHandler,
		authMiddleware,
		permMiddleware,
	)

	server := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router.SetupRoutes(),
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Failed to start server:", err)
	}
}
