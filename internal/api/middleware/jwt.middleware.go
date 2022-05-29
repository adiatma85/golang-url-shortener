package middleware

import (
	"net/http"
	"strings"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/repository"
	"github.com/adiatma85/golang-rest-template-api/pkg/crypto"
	"github.com/adiatma85/golang-rest-template-api/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Func to authorizing jwt token
func AuthJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := response.BuildFailedResponse("no token provided", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		requestToken := strings.Split(authHeader, " ")[1]
		jwtHelper := crypto.GetJWTCrypto()
		token, isValid, err := jwtHelper.ValidateToken(requestToken)
		if !isValid {
			response := response.BuildFailedResponse("token is not valid", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		claim, _ := jwtHelper.ExtractClaim(token)
		var smapClaim crypto.JwtCustomClaim
		mapstructure.Decode(claim, &smapClaim)

		// Get User Repository and set it to gin context
		userRepo := repository.GetUserRepository()
		user, _ := userRepo.GetById(smapClaim.UserID)
		c.Set("user", user)

		c.Next()
	}
}
