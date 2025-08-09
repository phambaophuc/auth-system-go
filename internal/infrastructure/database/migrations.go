package database

import (
	"auth-system/internal/domain/entities"
	"log"

	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&entities.User{},
		&entities.Role{},
		&entities.Permission{},
	); err != nil {
		return err
	}

	// Seed default data
	if err := seedDefaultData(db); err != nil {
		return err
	}

	return nil
}

func seedDefaultData(db *gorm.DB) error {
	// Create default permissions
	permissions := []entities.Permission{
		{Name: "users.read", Resource: "users", Action: "read", Description: "Read user information"},
		{Name: "users.write", Resource: "users", Action: "write", Description: "Create and update users"},
		{Name: "users.delete", Resource: "users", Action: "delete", Description: "Delete users"},
		{Name: "roles.read", Resource: "roles", Action: "read", Description: "Read role information"},
		{Name: "roles.write", Resource: "roles", Action: "write", Description: "Create and update roles"},
		{Name: "roles.delete", Resource: "roles", Action: "delete", Description: "Delete roles"},
	}

	for _, perm := range permissions {
		var existingPerm entities.Permission
		if err := db.Where("name = ?", perm.Name).First(&existingPerm).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&perm).Error; err != nil {
					log.Printf("Error creating permission %s: %v", perm.Name, err)
				}
			}
		}
	}

	// Create default roles
	roles := []struct {
		role        entities.Role
		permissions []string
	}{
		{
			role: entities.Role{
				Name:        "admin",
				Description: "Administrator with full access",
			},
			permissions: []string{"users.read", "users.write", "users.delete", "roles.read", "roles.write", "roles.delete"},
		},
		{
			role: entities.Role{
				Name:        "user",
				Description: "Regular user with limited access",
			},
			permissions: []string{"users.read"},
		},
	}

	for _, roleData := range roles {
		var existingRole entities.Role
		if err := db.Where("name = ?", roleData.role.Name).First(&existingRole).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create role
				if err := db.Create(&roleData.role).Error; err != nil {
					log.Printf("Error creating role %s: %v", roleData.role.Name, err)
					continue
				}

				// Assign permissions to role
				var perms []entities.Permission
				if err := db.Where("name IN ?", roleData.permissions).Find(&perms).Error; err == nil {
					if err := db.Model(&roleData.role).Association("Permissions").Append(&perms); err != nil {
						log.Printf("Error assigning permissions to role %s: %v", roleData.role.Name, err)
					}
				}
			}
		}
	}

	return nil
}
