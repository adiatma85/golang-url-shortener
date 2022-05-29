package repository

import (
	"strconv"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/db"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
)

// Local Variable

// Can save
// Can laod (multiple)
// Can load (single)
// Can Edit
// Can delete

// Check whether it's exist or not from shorten url

// Url that need to be returned when calling this instance of repository
var urlRepository *UrlRepository

// Interface that define what function that exist within struc
type UrlRepositoryInterface interface {
	Create(url models.Url) (models.Url, error)
	Query(pagination helpers.Pagination) (*helpers.Pagination, error)
	GetByShortenUniqueId(uniqueId string) (*models.Url, error)
	GetById(urlId string) (*models.Url, error)
	Update(url *models.Url) error
	Delete(url *models.Url) error

	// Add on or unique case
	// IsExist(uniqueId string) (bool, error)
}

// Struct that need to implement #UrlRepositoryInterface
type UrlRepository struct {
}

// Func to return Url Repository Instance
func GetUrlRepository() UrlRepositoryInterface {
	if urlRepository == nil {
		urlRepository = &UrlRepository{}
	}
	return urlRepository
}

// Func to Crreate Url and save it to database
func (repo *UrlRepository) Create(url models.Url) (models.Url, error) {
	err := Create(&url)
	// If error when transaction to database i.e duplicate email
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

// Func to Query of WHERE with pagination
func (repo *UrlRepository) Query(pagination helpers.Pagination) (*helpers.Pagination, error) {
	var urls []models.Url
	outputPagination, err := Query(&models.Url{}, &urls, pagination, []string{"User"})
	// If error when searching for query
	if err != nil {
		return nil, err
	}
	return outputPagination, nil
}

// Func to return model from random generated unique id
func (repo *UrlRepository) GetByShortenUniqueId(uniqueId string) (*models.Url, error) {
	var url models.Url
	where := models.Url{}
	where.ShortenUrl = uniqueId
	_, err := First(&where, &url, []string{"User"})
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// Func to GetById
func (repo *UrlRepository) GetById(urlId string) (*models.Url, error) {
	var url models.Url
	where := models.Url{}
	where.ID, _ = strconv.ParseUint(urlId, 10, 64)
	_, err := First(&where, &url, []string{"User"})
	if err != nil {
		return nil, err
	}
	return &url, nil
}

// Func to Update particular url
func (repo *UrlRepository) Update(url *models.Url) error {
	return Save(url)
}

// Func to Delete particular url
func (repo *UrlRepository) Delete(url *models.Url) error {
	_, err := DeleteByModel(url)
	if err != nil {
		return err
	}
	return nil
}

// Add on to check whether the random generated id is already exist or not
func (repo *UrlRepository) IsExist(uniqueId string) (bool, error) {
	// Own solution
	// isExistItem, err := repo.GetByShortenUniqueId(uniqueId)
	// if err != nil {
	// 	return false, err
	// }

	// return isExistItem != nil, nil

	// Another solution --> https://stackoverflow.com/questions/66392372/select-exists-with-gorm
	var exists bool
	database := db.GetDB()
	err := database.Model(models.Url{}).Where("shorten_url = ?", uniqueId).Find(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}
