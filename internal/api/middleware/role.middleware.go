package middleware

import (
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/pkg/crypto"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Is AdminMiddleware
func IsAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaim := c.MustGet("user_claim")
		var smapClaim crypto.JwtCustomClaim
		mapstructure.Decode(userClaim, &smapClaim)
		userRepo := repository.GetUserRepository()

		user, _ := userRepo.GetById(smapClaim.UserID)
		if user.Role.Name != "ADMIN" {
			response := response.BuildFailedResponse("you do not have permission to access this request", "unauhtorized request")
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}

		c.Next()
	}
}
