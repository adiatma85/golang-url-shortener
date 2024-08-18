package entity

import (
	"github.com/adiatma85/own-go-sdk/null"
	"github.com/adiatma85/own-go-sdk/query"
)

const (
	Base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Redis Key For counting
	UrlCountingRedisKey = "url_shortener:key_count:%s"
)

type Url struct {
	ID          int64       `db:"id" json:"id"`
	UserId      int64       `db:"fk_user_id" json:"userId"`
	OriginalUrl string      `db:"original_url" json:"originalUrl"`
	ShortenUrl  string      `db:"shorten_url" json:"shorten_url"`
	Visit       int64       `db:"visit" json:"visit"`
	Status      null.Int64  `db:"status" json:"status" swaggertype:"integer"`
	CreatedAt   null.Time   `db:"created_at" json:"createdAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	CreatedBy   null.String `db:"created_by" json:"createdBy" swaggertype:"string"`
	UpdatedAt   null.Time   `db:"updated_at" json:"updatedAt" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy   null.String `db:"updated_by" json:"updatedBy" swaggertype:"string"`
	DeletedAt   null.Time   `db:"deleted_at" json:"deletedAt,omitempty" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy   null.String `db:"deleted_by" json:"deletedBy,omitempty" swaggertype:"string"`
}

type UrlParam struct {
	ID         null.Int64 `param:"id" db:"id" uri:"url_id" form:"id"`
	UserId     null.Int64 `param:"fk_user_id" db:"fk_user_id" uri:"user_id" form:"fk_user_id"`
	ShortenUrl string     `param:"shorten_url" uri:"shorten_url" db:"shorten_url" swaggertype:"string"`
	PaginationParam
	QueryOption query.Option
}

type CreateUrlParam struct {
	UserId      int64       `db:"fk_user_id" json:"-"`
	OriginalUrl string      `db:"original_url" json:"originalUrl"`
	ShortenUrl  string      `db:"shorten_url" json:"-"`
	Visit       int64       `db:"visit" json:"-"`
	CreatedBy   null.String `json:"-" db:"created_by" swaggertype:"string"`
	UpdatedBy   null.String `json:"-" db:"updated_by" swaggertype:"string"`
}

type UpdateUrlParam struct {
	Visit     int64       `db:"visit" json:"visit"`
	Status    null.Int64  `param:"status" db:"status" json:"-" swaggertype:"integer"`
	UpdatedAt null.Time   `param:"updated_at" db:"updated_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	UpdatedBy null.String `param:"updated_by" db:"updated_by" json:"-" swaggertype:"string"`
	DeletedAt null.Time   `param:"deleted_at" db:"deleted_at" json:"-" swaggertype:"string" example:"2022-06-21T10:32:29Z"`
	DeletedBy null.String `param:"deleted_by" db:"deleted_by" json:"-" swaggertype:"string"`
}
