package tokenmanager

import (
	"fmt"
)

type ErrorInvalidClaims struct {
	msg string
}

func (e *ErrorInvalidClaims) Error() string {
	return fmt.Sprintf("invalid claims: %s", e.msg)
}

type ErrorParseDuration struct {
	msg string
}

func (e *ErrorParseDuration) Error() string {
	return fmt.Sprintf("parse duration: %s", e.msg)
}

type ErrorParseToken struct {
	msg string
}

func (e *ErrorParseToken) Error() string {
	return fmt.Sprintf("parse token: %s", e.msg)
}

type ErrorTokenSigningMethodMismatch struct {
	msg string
}

func (e *ErrorTokenSigningMethodMismatch) Error() string {
	return fmt.Sprintf("token signing method mismatch: %s", e.msg)
}

type ErrorUnsupportedSigningMethod struct {
	msg string
}

func (e *ErrorUnsupportedSigningMethod) Error() string {
	return fmt.Sprintf("unsupported signing method: %s", e.msg)
}
