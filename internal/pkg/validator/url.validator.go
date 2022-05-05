package validator

// Struct that define the validator/binding of Create Url Request
type CreateUrlRequest struct {
	OriginalUrl string `json:"original_url" binding:"required"`
}
