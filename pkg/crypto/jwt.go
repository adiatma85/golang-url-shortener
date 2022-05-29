package crypto

import (
	"fmt"
	"time"

	"github.com/adiatma85/golang-rest-template-api/internal/pkg/config"
	"github.com/golang-jwt/jwt"
)

var jwtHelper *jwtCryptoHelper

// Contract fot JWT Crypto Helper
type JWTCryptoHelper interface {
	GenerateToken(UserId string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, bool, error)
	ExtractClaim(token *jwt.Token) (jwt.MapClaims, bool)
}

// Struct for jwt custom claim
type JwtCustomClaim struct {
	UserID string `json:"user_id" mapstructure:"user_id"`
	jwt.StandardClaims
}

// Struct for JWTHelper
type jwtCryptoHelper struct {
}

// Func to initialize new jwt crypto helper
func GetJWTCrypto() JWTCryptoHelper {
	if jwtHelper == nil {
		jwtHelper = &jwtCryptoHelper{}
	}
	return jwtHelper
}

// Func to Generate Token with User ID as main issuer
func (helper *jwtCryptoHelper) GenerateToken(UserID string) (string, error) {
	serverConfiguration := config.GetConfig().Server
	claims := &JwtCustomClaim{
		UserID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(serverConfiguration.ExpiresHour)).Unix(),
			Issuer:    serverConfiguration.Name,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(serverConfiguration.Secret))
	if err != nil {
		return err.Error(), err
	}
	return t, nil
}

// Func to validate token
func (helper *jwtCryptoHelper) ValidateToken(tokenString string) (*jwt.Token, bool, error) {
	serverConfiguration := config.GetConfig().Server
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there was an error")
		}
		return []byte(serverConfiguration.Secret), nil
	})
	if err != nil {
		return nil, false, err
	}
	return token, token.Valid, nil
}

// Func to extract claim
func (helper *jwtCryptoHelper) ExtractClaim(token *jwt.Token) (jwt.MapClaims, bool) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, ok
	}
	return nil, false
}
