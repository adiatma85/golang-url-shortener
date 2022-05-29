package middleware

import (
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/constant"
	"github.com/adiatma85/golang-rest-template-api/pkg/helpers"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
)

// Is AdminMiddleware
func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.ExtractUserFromClaim(c)
		if user.Role.Name == constant.ADMINROLE {
			c.Next()
			return
		}
		response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}
}

// IsUserMiddleware
func IsUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.ExtractUserFromClaim(c)
		if user.Role.Name != constant.USERROLE {
			c.Next()
			return
		}
		response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}
}

// IsAdminOrUserMiddleware
func IsAdminOrUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := helpers.ExtractUserFromClaim(c)
		if user.Role.Name == constant.ADMINROLE || user.Role.Name == constant.USERROLE {
			c.Next()
			return
		}
		response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
		c.AbortWithStatusJSON(http.StatusUnauthorized, response)
	}
}
