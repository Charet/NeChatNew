package models

import "github.com/golang-jwt/jwt"

type MyClaim struct {
	UserID   int
	Username string
	jwt.StandardClaims
}
