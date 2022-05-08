package handler

import "github.com/gin-gonic/gin"

// This will be Auth Handler

// Local Variable
var authHandler *AuthHandler

// Struct to implement contract of AuthInterface
type AuthHandler struct{}

// Contract of Auth Interface
type AuthHandlerInterface interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

// Func to return Auth Handler instance
func GetAuthHandler() AuthHandlerInterface {
	if authHandler == nil {
		authHandler = &AuthHandler{}
	}
	return authHandler
}

// Login Func
func (handler *AuthHandler) Login(c *gin.Context) {
	// Need to be implemented here
}

// Register Func
func (handler *AuthHandler) Register(c *gin.Context) {
	// Need to be implemented here
}
