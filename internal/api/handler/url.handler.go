package handler

import (
	"net/http"
	"strconv"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/validator"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
)

// In here we declare url handler
var urlHandler *UrlHandler

// Struct that need to be implemented according to UrlHandlerInterface
type UrlHandler struct{}

// Interface contract for this instance
type UrlHandlerInterface interface {
	Create(c *gin.Context)
	Query(c *gin.Context)
	Load(c *gin.Context)
	// Function that need to be implemented
	// AuthorizedCreate(c *gin.Context)
	// AuthorizedDelete(c *gin.Context)
}

// Func to get instance of url handler
func GetUrlHandler() UrlHandlerInterface {
	if urlHandler == nil {
		urlHandler = &UrlHandler{}
	}
	return urlHandler
}

// Func to handle the request about new url
func (handler *UrlHandler) Create(c *gin.Context) {
	var newUrlRequest validator.CreateUrlRequest
	err := c.ShouldBind(&newUrlRequest)

	// Bad request error
	if err != nil {
		response := response.BuildFailedResponse("failed to add new shortener", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Get Url repo
	urlRepo := repository.GetUrlRepository()
	// Generate new randomized token
	urlModel := &models.Url{}
	smapping.FillStruct(urlModel, smapping.MapFields(&newUrlRequest))

	// For now it's some kind of hard-coded
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

// Func to handler the request to query urls
func (handler *UrlHandler) Query(c *gin.Context) {
	pagination := helpers.Pagination{}
	urlRepo := repository.GetUrlRepository()
	queryPageLimit, isPageLimitExist := c.GetQuery("limit")
	queryPage, isPageQueryExist := c.GetQuery("page")

	if isPageQueryExist {
		pagination.Page, _ = strconv.Atoi(queryPage)
	}

	if isPageLimitExist {
		pagination.Limit, _ = strconv.Atoi(queryPageLimit)
	}

	urls, err := urlRepo.Query(pagination)

	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	response := response.BuildSuccessResponse("success to fetch data", urls)
	c.JSON(http.StatusOK, response)
}

// Func to return the info of the instance of url with given shorten version of url
func (handler *UrlHandler) Load(c *gin.Context) {
	shortToken := c.Param("short_token")
	// Get Url repo
	urlRepo := repository.GetUrlRepository()

	url, err := urlRepo.GetByShortenUniqueId(shortToken)

	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := response.BuildSuccessResponse("success to fetch data", url)
	c.JSON(http.StatusOK, response)
}

// Add user entity, repository, and handler
