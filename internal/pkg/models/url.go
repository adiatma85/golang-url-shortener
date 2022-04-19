package models

// Struct for Url Models
type Url struct {
	Model
	OriginalUrl string `gorm:"type:varchar(255)" json:"url"`
	ShortenUrl  string `gorm:"type:varchar(255)" json:"shorten_url"`
}
