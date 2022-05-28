package validator

// Struct that define the validator/binding of Create Url Request
type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}

// Struct that define the validator/binding of Authorized Url Request
type AuthorizedCreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortenUrl  string `json:"shorten_url" binding:"required"`
}
