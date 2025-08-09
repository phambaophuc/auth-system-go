package services

import (
	"auth-system/internal/domain/entities"
	"auth-system/internal/domain/repositories"
	"auth-system/internal/domain/services"
	"fmt"
)

type permissionService struct {
	permissionRepo repositories.PermissionRepository
	userRepo       repositories.UserRepository
}

func NewPermissionService(
	permissionRepo repositories.PermissionRepository,
	userRepo repositories.UserRepository,
) services.PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		userRepo:       userRepo,
	}
}

func (s *permissionService) CheckPermission(userID uint, resource, action string) (bool, error) {
	permissions, err := s.permissionRepo.GetByUserID(userID)
	if err != nil {
		return false, fmt.Errorf("Failed to get user permissions: %w", err)
	}

	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
	}

	return false, nil
}

func (s *permissionService) GetUserPermissions(userID uint) ([]*entities.Permission, error) {
	return s.permissionRepo.GetByUserID(userID)
}
