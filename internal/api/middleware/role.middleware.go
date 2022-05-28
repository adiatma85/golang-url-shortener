package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Is AdminMiddleware
func IsAdminMiddleware(c *gin.Context) {
	// for now testing for return claim
	c.JSON(http.StatusOK, c.MustGet("user_claim"))
}
