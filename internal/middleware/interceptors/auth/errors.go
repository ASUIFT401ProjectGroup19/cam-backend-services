package auth

import (
	"fmt"
)

type ErrorAuthHeader struct {
	msg string
}

func (e *ErrorAuthHeader) Error() string {
	return fmt.Sprintf("auth header: %s", e.msg)
}

type ErrorMetadata struct {
	msg string
}

func (e *ErrorMetadata) Error() string {
	return fmt.Sprintf("metadata: %s", e.msg)
}

type ErrorTokenValidation struct {
	msg string
}

func (e *ErrorTokenValidation) Error() string {
	return fmt.Sprintf("token validation: %s", e.msg)
}
