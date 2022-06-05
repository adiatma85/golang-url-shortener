package helpers

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
)

// Reference --> https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
// Function below used to create random fixed string. It used to make shortener version of link
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrcSB(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}

// Helper Function below used in Role Middleware and Handler to extract User Entity from
func ExtractUserFromClaim(c *gin.Context) *models.User {
	user := c.MustGet("user")
	var exstingUser models.User
	mapstructure.Decode(user, &exstingUser)
	return &exstingUser
}

// Helper to change string to uint
func ConvertStringtoUint(id string) uint {
	typeUint64, _ := strconv.ParseUint(id, 10, 64)
	return uint(typeUint64)
}
