package middleware

import (
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/constant"
	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Is AdminMiddleware
func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := extractUserFromClaim(c)
		if user.Role.Name != constant.ADMINROLE {
			response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
		c.Next()
	}
}

// IsUserMiddleware
func IsUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := extractUserFromClaim(c)
		if user.Role.Name != constant.USERROLE {
			response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}
		c.Next()
	}
}

// Helper function to extract user
func extractUserFromClaim(c *gin.Context) *models.User {
	user := c.MustGet("user")
	var exstingUser models.User
	mapstructure.Decode(user, &exstingUser)
	return &exstingUser
}
