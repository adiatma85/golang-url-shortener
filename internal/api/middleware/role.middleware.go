package middleware

import (
	"fmt"
	"net/http"

	"github.com/adiatma85/golang-rest-template-api/pkg/crypto"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Is AdminMiddleware
func IsAdminMiddleware(c *gin.Context) {
	// for now testing for return claim
	userClaim := c.MustGet("user_claim")
	var smapClaim crypto.JwtCustomClaim
	mapstructure.Decode(userClaim, &smapClaim)

	fmt.Println(userClaim)
	c.JSON(http.StatusOK, smapClaim.UserID)
}
