package services

import (
	"auth-system/internal/application/dto"
	"auth-system/internal/domain/entities"
)

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, error)
	Register(req *dto.RegisterRequest) (*dto.AuthResponse, error)
	RefreshToken(refreshToken string) (*dto.AuthResponse, error)
	Logout(userID uint) error
}

type UserService interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
	AssignRole(userID uint, roleID uint) error
	RemoveRole(userID uint, roleID uint) error
}

type PermissionService interface {
	CheckPermission(userID uint, resource, action string) (bool, error)
	GetUserPermissions(userID uint) ([]*entities.Permission, error)
}
