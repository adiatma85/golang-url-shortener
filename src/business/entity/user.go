package entity

import (
	"github.com/adiatma85/own-go-sdk/jwtAuth"
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
)

type User struct {
	ID          int64       `db:"id" json:"id"`
	RoleId      null.Int64  `db:"fk_role_id" json:"roleId"`
	Email       string      `db:"email" json:"email"`
	Username    string      `db:"username" json:"username"`
	Password    string      `db:"password" json:"-"`
	DisplayName string      `db:"display_name" json:"displayName"`
	Status      null.Int64  `db:"status" json:"status" swaggertype:"integer"`
	CreatedAt   null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy   null.String `db:"created_by" json:"createdBy" swaggertype:"string"`
	UpdatedAt   null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy   null.String `db:"updated_by" json:"updatedBy" swaggertype:"string"`
	DeletedAt   null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy   null.String `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"string"`
}

func (u *User) ConvertToAuthUser() jwtAuth.User {
	return jwtAuth.User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
	}
}

type UserParam struct {
	ID          null.Int64  `param:"id" uri:"user_id" db:"id" form:"id"`
	RoleId      null.Int64  `param:"fk_role_id" uri:"role_id" db:"fk_role_id" form:"fk_role_id"`
	IDs         []int64     `param:"ids" uri:"user_ids" db:"id" form:"userIds"`
	Email       null.String `param:"email" db:"email"`
	Username    null.String `param:"username" db:"username"`
	DisplayName null.String `param:"display_name" db:"display_name"`
	PaginationParam
	QueryOption query.Option
}

type CreateUserParam struct {
	RoleId          int64       `db:"fk_role_id" json:"-"`
	Email           string      `db:"email" json:"email"`
	Username        string      `db:"username" json:"username"`
	Password        string      `db:"password" json:"password"`
	ConfirmPassword string      `db:"-" json:"confirmPassword"`
	DisplayName     string      `db:"display_name" json:"displayName"`
	CreatedBy       null.String `json:"-" db:"created_by" swaggertype:"string"`
	UpdatedBy       null.String `json:"-" db:"updated_by" swaggertype:"string"`
}

type UpdateUserParam struct {
	RoleId      string      `param:"fk_role_id" db:"fk_role_id" json:"roleId"`
	Username    string      `param:"username" db:"username" json:"username"`
	DisplayName string      `param:"display_name" db:"display_name" json:"displayName"`
	Password    string      `param:"password" db:"password" json:"-"`
	Status      null.Int64  `param:"status" db:"status" json:"-" swaggertype:"integer"`
	UpdatedAt   null.Time   `param:"updated_at" db:"updated_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy   null.String `param:"updated_by" db:"updated_by" json:"-" swaggertype:"string"`
	DeletedAt   null.Time   `param:"deleted_at" db:"deleted_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy   null.String `param:"deleted_by" db:"deleted_by" json:"-" swaggertype:"string"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	AccessToken string `json:"accessToken"`
}

type ChangePasswordRequest struct {
	OldPassword     string `db:"-" json:"oldPassword"`
	Password        string `db:"-" json:"newPassword"`
	ConfirmPassword string `db:"-" json:"confirmPassword"`
}
