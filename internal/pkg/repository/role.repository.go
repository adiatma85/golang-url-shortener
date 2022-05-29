package repository

import "github.com/adiatma85/golang-rest-template-api/internal/pkg/models"

// Local variable
var (
	roleRepo *RoleRepository
)

// Contract of Role Repository
type RoleRepositoryInterface interface {
	GetRoleByName(name string) (*models.Role, error)
}

// Struct ti implements contract or interface
type RoleRepository struct {
}

// Func to return Role Repository instance
func GetRoleRepository() RoleRepositoryInterface {
	if roleRepo == nil {
		roleRepo = &RoleRepository{}
	}
	return roleRepo
}

// Func to Get Role By Role Name
func (repo *RoleRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	where := models.Role{}
	where.Name = name
	_, err := First(&where, &role, []string{})
	if err != nil {
		return nil, err
	}
	return &role, nil
}
