package seeders

import (
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/constant"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
)

// Var Roles for seeder
var Roles []models.Role = []models.Role{
	{
		Name: constant.ADMINROLE,
	},
	{
		Name: constant.USERROLE,
	},
}
