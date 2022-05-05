package handler

import (
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/validator"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

// Can Create
// Can get from id

// In here we declare url handler
var urlHandler *UrlHandler

// Struct that need to be implemented according to UrlHandlerInterface
type UrlHandler struct{}

// Interface contract for this instance
type UrlHandlerInterface interface {
	UrlCreate(c *gin.Context)
	UrlLoad(c *gin.Context)
}

// Func to get instance of url handler
func GetUrlHandler() UrlHandlerInterface {
	if urlHandler == nil {
		urlHandler = &UrlHandler{}
	}
	return urlHandler
}

// func to handle the request about new url
func (handler *UrlHandler) UrlCreate(c *gin.Context) {
	var newUrlRequest validator.CreateUrlRequest
	err := c.ShouldBind(&newUrlRequest)

	// Bad request error
	if err != nil {
		response := response.BuildFailedResponse("failed to add new shortener", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Get Url repo
	urlRepo := repository.GetUrlRepository()
	// Generate new randomized token
	urlModel := &models.Url{}
	smapping.FillStruct(urlModel, smapping.MapFields(&newUrlRequest))

	urlModel.ShortenUrl = helpers.RandStringBytesMaskImprSrcSB(6)

	if newUrl, err := urlRepo.Create(*urlModel); err != nil {
		response := response.BuildFailedResponse("failed to add new shortener", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	} else {
		data := map[string]interface{}{
			"original_url": newUrlRequest.OriginalUrl,
			"short_url":    newUrl.ShortenUrl,
		}
		response := response.BuildSuccessResponse("success to add new shortener url", data)
		c.JSON(http.StatusOK, response)
		return
	}
}

// Func to return the info of the instance of url with given shorten version of url
func (handler *UrlHandler) UrlLoad(c *gin.Context) {

}
