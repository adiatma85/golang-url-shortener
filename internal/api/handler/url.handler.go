package handler

import (
	"net/http"
	"strconv"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/config"
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
type UrlHandler struct {
	UrlCharacterLong int
	UrlRepo          repository.UrlRepositoryInterface
}

// Interface contract for this instance
type UrlHandlerInterface interface {
	Create(c *gin.Context)
	Query(c *gin.Context)
	Load(c *gin.Context)
	// Function that need to be implemented
	AuthorizedCreate(c *gin.Context)
	AuthorizedUpdate(c *gin.Context)
	AuthorizedDelete(c *gin.Context)
}

// Func to get instance of url handler
func GetUrlHandler() UrlHandlerInterface {
	if urlHandler == nil {
		urlHandler = &UrlHandler{
			UrlCharacterLong: config.GetConfig().Server.UrlCharacterLong,
			UrlRepo:          repository.GetUrlRepository(),
		}
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
	urlRepo := handler.UrlRepo
	// Generate new randomized token
	urlModel := &models.Url{}
	smapping.FillStruct(urlModel, smapping.MapFields(&newUrlRequest))

	// For now it's some kind of hard-coded
	user := helpers.ExtractUserFromClaim(c)
	urlModel.ShortenUrl = helpers.RandStringBytesMaskImprSrcSB(handler.UrlCharacterLong)
	urlModel.UserId = user.ID

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
	urlRepo := handler.UrlRepo
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
	urlRepo := handler.UrlRepo

	url, err := urlRepo.GetByShortenUniqueId(shortToken)

	if err != nil {
		response := response.BuildFailedResponse("failed to fetch data", err.Error())
		c.AbortWithStatusJSON(http.StatusNotFound, response)
		return
	}

	response := response.BuildSuccessResponse("success to fetch data", url)
	c.JSON(http.StatusOK, response)
}

// Func to authorized New Url when user can customize it
func (handler *UrlHandler) AuthorizedCreate(c *gin.Context) {
	var newUrlRequest validator.AuthorizedCreateUrlRequest
	err := c.ShouldBind(&newUrlRequest)

	// Bad request error
	if err != nil {
		response := response.BuildFailedResponse("failed to add new shortener", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	// Get Url Repo
	urlRepo := handler.UrlRepo
	// Generate new randomized token
	urlModel := &models.Url{}
	smapping.FillStruct(urlModel, smapping.MapFields(&newUrlRequest))

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

// Func to Update existed URL (and it's own of their respective user)
func (handler *UrlHandler) AuthorizedUpdate(c *gin.Context) {
	var updateRequest validator.AuthorizedUpdateRequest
	err := c.ShouldBind(&updateRequest)

	// Bad request error
	if err != nil {
		response := response.BuildFailedResponse("failed to add new shortener", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	updateModel := &models.Url{}

	// smapping the update request to models
	updateModel.ID = helpers.ConvertStringtoUint(c.Param("id"))
	smapping.FillStruct(updateModel, smapping.MapFields(&updateRequest))

	urlRepo := handler.UrlRepo
	err = urlRepo.Update(updateModel)

	if err != nil {
		response := response.BuildFailedResponse("failed to update a url", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Func to authorized Delete existed URL (and it's own of their respective user)
func (handler *UrlHandler) AuthorizedDelete(c *gin.Context) {
	deleteModel := &models.Url{}
	deleteModel.ID = helpers.ConvertStringtoUint(c.Param("id"))

	urlRepo := handler.UrlRepo

	err := urlRepo.Delete(deleteModel)
	if err != nil {
		response := response.BuildFailedResponse("failed to delete a url item", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
