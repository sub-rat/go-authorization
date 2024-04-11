package dto

import (
	"github.com/golang-jwt/jwt"
)

type JwtClaims struct {
	ID       string
	Username string
	jwt.StandardClaims
}
