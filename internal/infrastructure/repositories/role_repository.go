package repositories

import (
	"auth-system/internal/domain/entities"
	"auth-system/internal/domain/repositories"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repositories.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *entities.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uint) (*entities.Role, error) {
	var role entities.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*entities.Role, error) {
	var role entities.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *entities.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Role{}, id).Error
}

func (r *roleRepository) List() ([]*entities.Role, error) {
	var roles []*entities.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}
