package errs

import (
	"fmt"
)

type InvalidClaims struct {
	Message string
}

func (e *InvalidClaims) Error() string {
	return fmt.Sprintf("invalid claims: %s", e.Message)
}

type ParseDuration struct {
	Message string
}

func (e *ParseDuration) Error() string {
	return fmt.Sprintf("parse duration: %s", e.Message)
}

type ParseToken struct {
	Message string
}

func (e *ParseToken) Error() string {
	return fmt.Sprintf("parse token: %s", e.Message)
}

type TokenSigningMethodMismatch struct {
	Message string
}

func (e *TokenSigningMethodMismatch) Error() string {
	return fmt.Sprintf("token signing method mismatch: %s", e.Message)
}

type UnsupportedSigningMethod struct {
	Message string
}

func (e *UnsupportedSigningMethod) Error() string {
	return fmt.Sprintf("unsupported signing method: %s", e.Message)
}
