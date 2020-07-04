package common

import "github.com/dgrijalva/jwt-go"

type JWTClaims struct {
	UUID      string
	AccountID string
	jwt.StandardClaims
}
