package models

import (
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/config"
	"gorm.io/gorm"
)

// Struct for Url Models
type Url struct {
	Model
	OriginalUrl string `gorm:"type:varchar(255)" json:"url"`
	ShortenUrl  string `gorm:"type:varchar(255)" json:"shorten_url"`
}

// Hook Afterfind to return formatted url in shorten url
func (u *Url) AfterFind(tx *gorm.DB) (err error) {
	urlEndpoint := "/api/v1/url/load"
	u.ShortenUrl = config.GetConfig().Server.Endpoint + urlEndpoint + "/" + u.ShortenUrl
	return
}
