package handler

import (
	"fmt"
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/validator"
	"github.com/adiatma85/golang-rest-template-api/pkg/crypto"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mashingan/smapping"
	"github.com/mitchellh/mapstructure"
)

// Local Variable
var authHandler *AuthHandler

// Struct to implement contract of AuthInterface
type AuthHandler struct {
	UserRepository repository.UserRepositoryInterface
}

// Contract of Auth Interface
type AuthHandlerInterface interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
	DeleteProfile(c *gin.Context)
}

// Func to return Auth Handler instance
func GetAuthHandler() AuthHandlerInterface {
	if authHandler == nil {
		authHandler = &AuthHandler{
			UserRepository: repository.GetUserRepository(),
		}
	}
	return authHandler
}

// Login Func
func (handler *AuthHandler) Login(c *gin.Context) {
	var loginRequest validator.LoginRequest
	err := c.ShouldBind(&loginRequest)

	// Error when binding in validator
	if err != nil {
		response := response.BuildFailedResponse("failed to login", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userRepo := handler.UserRepository
	// If user doesn't exist
	if user, err := userRepo.GetByEmail(loginRequest.Email); err != nil {
		response := response.BuildFailedResponse("wrong credential", "username or password is not match with our database")
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	} else {
		// Wrong password
		passwordHelper := crypto.GetPasswordCryptoHelper()
		if !passwordHelper.ComparePassword(user.Password, []byte(loginRequest.Password)) {
			response := response.BuildFailedResponse("wrong credential", "username or password is not match with our database")
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}

		// Correct password
		tokenHelper := crypto.GetJWTCrypto()
		token, err := tokenHelper.GenerateToken(fmt.Sprint(user.ID))
		// Error when creating new token
		if err != nil {
			response := response.BuildFailedResponse("failed to generate token", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		response := response.BuildSuccessResponse("success login", map[string]interface{}{
			"token": token,
		})
		c.JSON(http.StatusOK, response)
		return
	}
}

// Register Func
func (handler *AuthHandler) Register(c *gin.Context) {
	var registerRequest validator.RegisterRequest
	err := c.ShouldBind(&registerRequest)

	// Error when binding in validator
	if err != nil {
		response := response.BuildFailedResponse("failed to register", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userRepo := handler.UserRepository
	passwordHelper := crypto.GetPasswordCryptoHelper()
	userModel := &models.User{}

	// smapping the struct
	smapping.FillStruct(userModel, smapping.MapFields(&registerRequest))
	userModel.Password, _ = passwordHelper.HashAndSalt([]byte(registerRequest.Password))

	if newUser, err := userRepo.Create(*userModel); err != nil {
		response := response.BuildFailedResponse("failed to register", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	} else {
		tokenHelper := crypto.GetJWTCrypto()
		token, err := tokenHelper.GenerateToken(fmt.Sprint(newUser.ID))
		if err != nil {
			response := response.BuildFailedResponse("failed to generate token", err.Error())
			c.AbortWithStatusJSON(http.StatusInternalServerError, response)
			return
		}
		response := response.BuildSuccessResponse("success register new user", map[string]interface{}{
			"token": token,
		})
		c.JSON(http.StatusOK, response)
		return
	}
}

// Func to Get each User Profile
func (handler *AuthHandler) GetProfile(c *gin.Context) {
	user := helpers.ExtractUserFromClaim(c)
	c.JSON(http.StatusOK, user)
}

// UpdateProfile Func
func (handler *AuthHandler) UpdateProfile(c *gin.Context) {
	var updateProfileRequest validator.UpdateProfileRequest
	err := c.ShouldBind(&updateProfileRequest)

	// Error when binding in validator
	if err != nil {
		response := response.BuildFailedResponse("failed to register", err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, response)
		return
	}

	userRepo := handler.UserRepository

	// Existing user from middleware
	existedUser := c.MustGet("user")
	var mappedUser models.User
	mapstructure.Decode(existedUser, &mappedUser)

	smapping.FillStruct(&mappedUser, smapping.MapFields(&updateProfileRequest))

	err = userRepo.Update(&mappedUser)

	if err != nil {
		response := response.BuildFailedResponse("failed to update profile", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Delete Profile Func
func (handler *AuthHandler) DeleteProfile(c *gin.Context) {
	userRepo := handler.UserRepository
	// Existing user from middleware
	existedUser := c.MustGet("user")
	var mappedUser models.User
	mapstructure.Decode(existedUser, &mappedUser)

	err := userRepo.Delete(&mappedUser)

	if err != nil {
		response := response.BuildFailedResponse("failed to delete profile", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, response)
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
