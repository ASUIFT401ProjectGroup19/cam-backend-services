package auth

import (
	"fmt"
)

type MissingHeader struct {
	Message string
}

func (e *MissingHeader) Error() string {
	return fmt.Sprintf("missing header: %s", e.Message)
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
