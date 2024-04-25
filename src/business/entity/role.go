package entity

import (
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
)

const (
	RoleIdSuperAdmin = 1

	// Role Type Enum
	RoleTypeAdmin = "admin"
	RoleTypeUser  = "user"
)

type Role struct {
	ID        int64       `db:"id" json:"id"`
	Name      string      `db:"name" json:"name"`
	Type      string      `db:"type" json:"type"`
	Rank      int64       `db:"rank" json:"rank"`
	Status    int64       `db:"status" json:"status" swaggertype:"integer"`
	CreatedAt null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy null.String `db:"created_by" json:"createdBy" swaggertype:"string"`
	UpdatedAt null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy null.String `db:"updated_by" json:"updatedBy" swaggertype:"string"`
	DeletedAt null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy null.String `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"string"`
}

type RoleParam struct {
	ID     null.Int64  `param:"id" uri:"role_id" db:"id" form:"role_id"`
	IDs    []int64     `param:"ids" uri:"role_ids" db:"id"`
	Name   null.String `param:"name" uri:"role_name" db:"id" form:"role_name"`
	Type   null.String `param:"type" uri:"role_type" db:"id" form:"role_type"`
	Status null.Int64  `param:"status" db:"status" swaggertype:"string"`
	PaginationParam
	QueryOption query.Option
}

type CreateRoleParam struct {
	Name      string      `db:"name" json:"name"`
	Type      string      `db:"type" json:"type"`
	Rank      int64       `db:"rank" json:"rank"`
	CreatedBy null.String `json:"-" db:"created_by" swaggertype:"string"`
	UpdatedBy null.String `json:"-" db:"updated_by" swaggertype:"string"`
}

type UpdateRoleParam struct {
	Name      string      `db:"name" json:"name"`
	Type      string      `db:"type" json:"type"`
	Rank      int64       `db:"rank" json:"rank"`
	Status    null.Int64  `db:"status" param:"status" json:"status" swaggertype:"string"`
	UpdatedAt null.Time   `db:"updated_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy null.String `db:"updated_by" json:"-" swaggertype:"string"`
	DeletedAt null.Time   `db:"deleted_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy null.String `db:"deleted_by" json:"-" swaggertype:"string"`
}
