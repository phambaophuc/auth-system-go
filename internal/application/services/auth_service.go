package services

import (
	"auth-system/internal/application/dto"
	"auth-system/internal/domain/entities"
	"auth-system/internal/domain/repositories"
	"auth-system/internal/domain/services"
	"auth-system/internal/infrastructure/security"
	"auth-system/pkg/errors"
	"fmt"

	"gorm.io/gorm"
)

type authService struct {
	userRepo        repositories.UserRepository
	roleRepo        repositories.RoleRepository
	jwtManager      *security.JWTManager
	passwordManager *security.PasswordManager
}

func NewAuthService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	jwtManager *security.JWTManager,
	passwordManager *security.PasswordManager,
) services.AuthService {
	return &authService{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		jwtManager:      jwtManager,
		passwordManager: passwordManager,
	}
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewValidationError("Invalid credentials")
		}
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, errors.NewValidationError("Account is deactivated")
	}

	if err := s.passwordManager.CheckPassword(user.Password, req.Password); err != nil {
		return nil, errors.NewValidationError("Invalid credentials")
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.mapUserToResponse(user),
	}, nil
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Check if user already exists
	_, err := s.userRepo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.NewValidationError("Email already exists")
	}

	hashedPassword, err := s.passwordManager.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash password: %w", err)
	}

	user := &entities.User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("Failed to create user: %w", err)
	}

	// Assign default user role
	userRole, err := s.roleRepo.GetByName("user")
	if err == nil {
		user.Roles = append(user.Roles, *userRole)
		s.userRepo.Update(user)
	}

	accessToken, refreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         s.mapUserToResponse(user),
	}, nil
}

func (s *authService) RefreshToken(refreshToken string) (*dto.AuthResponse, error) {
	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.NewValidationError("Invalid refresh token")
	}

	if claims.Type != "refresh" {
		return nil, errors.NewValidationError("Invalid token type")
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	if !user.IsActive {
		return nil, errors.NewValidationError("Account is deactivated")
	}

	accessToken, newRefreshToken, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate tokens: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         s.mapUserToResponse(user),
	}, nil
}

func (s *authService) Logout(userID uint) error {
	// In a production system, you might want to blacklist the tokens
	// For now, we just return success
	return nil
}

func (s *authService) mapUserToResponse(user *entities.User) dto.UserResponse {
	roles := make([]dto.RoleResponse, len(user.Roles))
	for i, role := range user.Roles {
		permissions := make([]dto.PermissionResponse, len(role.Permissions))
		for j, perm := range role.Permissions {
			permissions[j] = dto.PermissionResponse{
				ID:          perm.ID,
				Name:        perm.Name,
				Resource:    perm.Resource,
				Action:      perm.Action,
				Description: perm.Description,
			}
		}
		roles[i] = dto.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissions,
		}
	}

	return dto.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		Roles:     roles,
	}
}
