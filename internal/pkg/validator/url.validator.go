package validator

// Struct that define the validator/binding of Create Url Request
type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}

// Struct that define the validator/binding of Authorized Create Url Request
type AuthorizedCreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
	ShortenUrl  string `json:"shorten_url" binding:"required"`
}

// Struct that define the validator/binding of Authorization Update Url Request
type AuthorizedUpdateRequest struct {
	OriginalUrl string `json:"original_url"`
	ShortenUrl  string `json:"shorten_url"`
}
