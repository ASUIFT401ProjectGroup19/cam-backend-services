package errs

import (
	"fmt"
)

type AuthHeader struct {
	Message string
}

func (e *AuthHeader) Error() string {
	return fmt.Sprintf("auth header: %s", e.Message)
}

type Metadata struct {
	Message string
}

func (e *Metadata) Error() string {
	return fmt.Sprintf("metadata: %s", e.Message)
}

type TokenValidation struct {
	Message string
}

func (e *TokenValidation) Error() string {
	return fmt.Sprintf("token validation: %s", e.Message)
}
