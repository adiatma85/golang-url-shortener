package models

import (
	"time"

	"gorm.io/gorm"
)

// Struct that derivative from gorm.Model but customized
type Model struct {
	ID        uint           `gorm:"column:id;primary_key;auto_increment;" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}
