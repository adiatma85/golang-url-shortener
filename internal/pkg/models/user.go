package models

// Struct for User Models
type User struct {
	Model
	Name     string `gorm:"type:varchar(100)" json:"name" validation:"name"`
	Email    string `gorm:"type:varchar(100);unique;" json:"email" validation:"email"`
	Password string `gorm:"type:varchar(100)" json:"-" validation:"password"`
}
