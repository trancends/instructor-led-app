package sharedmodel

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	jwt.RegisteredClaims
	AuthorID string
	Role     string
}
