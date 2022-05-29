package models

import "gorm.io/gorm"

// Struct for User Models
type User struct {
	Model
	Name     string `gorm:"type:varchar(100)" json:"name" validation:"name"`
	Email    string `gorm:"type:varchar(100);unique;" json:"email" validation:"email"`
	Password string `gorm:"type:varchar(100)" json:"-" validation:"password"`
	RoleId   uint64 `gorm:"not null" json:"-" validation:"role_id"`
	Role     Role   `gorm:"foreignkey:RoleId;constraint:onUpdate:CASCADE,onDelete:CASCADE" json:"role"`
}

// BeforeCreate Hook
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.RoleId == 0 {
		// Default Role is 2 for USER ROLE
		u.RoleId = 2
	}
	return
}
