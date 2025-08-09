package services

import (
	"auth-system/internal/application/dto"
	"auth-system/internal/domain/repositories"
	"auth-system/internal/domain/services"
	"auth-system/internal/infrastructure/security"
	"auth-system/pkg/errors"
	"fmt"

	"gorm.io/gorm"
)

type userService struct {
	userRepo        repositories.UserRepository
	roleRepo        repositories.RoleRepository
	passwordManager *security.PasswordManager
	authService     *authService
}

func NewUserService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	passwordManager *security.PasswordManager,
) services.UserService {
	return &userService{
		userRepo:        userRepo,
		roleRepo:        roleRepo,
		passwordManager: passwordManager,
		authService:     &authService{userRepo: userRepo, roleRepo: roleRepo},
	}
}

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User not found")
		}
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	userResponse := s.authService.mapUserToResponse(user)
	return &userResponse, nil
}

func (s *userService) UpdateProfile(userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewNotFoundError("User not found")
		}
		return nil, fmt.Errorf("Failed to get user: %w", err)
	}

	user.FirstName = req.FirstName
	user.LastName = req.LastName

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("Failed to update user: %w", err)
	}

	userResponse := s.authService.mapUserToResponse(user)
	return &userResponse, nil
}

func (s *userService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.NewNotFoundError("User not found")
		}
		return fmt.Errorf("Failed to get user: %w", err)
	}

	if err := s.passwordManager.CheckPassword(user.Password, req.CurrentPassword); err != nil {
		return errors.NewValidationError("Current password is incorrect")
	}

	hashedPassword, err := s.passwordManager.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("Failed to hash password: %w", err)
	}

	user.Password = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return fmt.Errorf("Failed to update password: %w", err)
	}

	return nil
}

func (s *userService) AssignRole(userID uint, roleID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("Failed to get user: %w", err)
	}

	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("Failed to get role: %w", err)
	}

	// Check if user already has this role
	for _, userRole := range user.Roles {
		if userRole.ID == roleID {
			return errors.NewValidationError("User already has this role")
		}
	}

	user.Roles = append(user.Roles, *role)
	return s.userRepo.Update(user)
}

func (s *userService) RemoveRole(userID uint, roleID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("Failed to get user: %w", err)
	}

	// Find and remove the role
	for i, role := range user.Roles {
		if role.ID == roleID {
			user.Roles = append(user.Roles[:i], user.Roles[i+1:]...)
			return s.userRepo.Update(user)
		}
	}

	return errors.NewValidationError("User does not have this role")
}
