package models

import (
	"time"

	"gorm.io/gorm"
)

// Struct that derivative from gorm.Model but customized
type Model struct {
	ID        uint64         `gorm:"column:id;primary_key;auto_increment;" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime;not null;" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at;type:datetime" json:"deleted_at"`
}
