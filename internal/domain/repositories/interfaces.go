package repositories

import "auth-system/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByID(id uint) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*entities.User, error)
}

type RoleRepository interface {
	Create(role *entities.Role) error
	GetByID(id uint) (*entities.Role, error)
	GetByName(name string) (*entities.Role, error)
	Update(role *entities.Role) error
	Delete(id uint) error
	List() ([]*entities.Role, error)
}

type PermissionRepository interface {
	Create(permission *entities.Permission) error
	GetByID(id uint) (*entities.Permission, error)
	GetByName(name string) (*entities.Permission, error)
	Update(permission *entities.Permission) error
	Delete(id uint) error
	List() ([]*entities.Permission, error)
	GetByUserID(userID uint) ([]*entities.Permission, error)
}
