package models

// Struct that represent Role
type Role struct {
	Model
	Name string `gorm:"type:varchar(100)" json:"name" validation:"name"`
}
