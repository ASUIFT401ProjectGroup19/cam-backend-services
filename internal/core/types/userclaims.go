package types

import "github.com/golang-jwt/jwt/v4"

type UserClaims struct {
	jwt.RegisteredClaims
}
