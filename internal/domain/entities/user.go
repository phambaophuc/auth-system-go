package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"unique;not null" json:"email"`
	Password  string         `gorm:"not null" json:"-"`
	FirstName string         `gorm:"not null" json:"first_name"`
	LastName  string         `gorm:"not null" json:"last_name"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	Roles     []Role         `gorm:"many2many:user_roles;" json:"roles"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
