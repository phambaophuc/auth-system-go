package repositories

import (
	"auth-system/internal/domain/entities"
	"auth-system/internal/domain/repositories"

	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) repositories.PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *entities.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id uint) (*entities.Permission, error) {
	var permission entities.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) GetByName(name string) (*entities.Permission, error) {
	var permission entities.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

func (r *permissionRepository) Update(permission *entities.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&entities.Permission{}, id).Error
}

func (r *permissionRepository) List() ([]*entities.Permission, error) {
	var permissions []*entities.Permission
	err := r.db.Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) GetByUserID(userID uint) ([]*entities.Permission, error) {
	var permissions []*entities.Permission

	query := `
		SELECT DISTINCT p.* FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN roles r ON rp.role_id = r.id
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = ?
	`

	err := r.db.Raw(query, userID).Scan(&permissions).Error
	return permissions, err
}
