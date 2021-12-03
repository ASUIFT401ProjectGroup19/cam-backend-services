package tokenmanager

import (
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errors"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Config struct {
	SecretKey     string
	SigningMethod string
	ValidDuration string
}

type UserClaims struct {
	jwt.RegisteredClaims
}

type TokenManager struct {
	secretKey     string
	signingMethod jwt.SigningMethod
	validDuration time.Duration
}

func New(config *Config) (*TokenManager, error) {
	var method jwt.SigningMethod
	switch config.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	default:
		return nil, &errors.UnsupportedSigningMethod{Message: config.SigningMethod}
	}
	duration, err := time.ParseDuration(config.ValidDuration)
	if err != nil {
		return nil, &errors.ParseDuration{Message: err.Error()}
	}
	return &TokenManager{
		secretKey:     config.SecretKey,
		signingMethod: method,
		validDuration: duration,
	}, nil
}

func (t *TokenManager) Generate(user *models.User) (string, error) {
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.validDuration)),
		},
	}
	token := jwt.NewWithClaims(t.signingMethod, claims)
	return token.SignedString([]byte(t.secretKey))
}

func (t *TokenManager) Validate(rawToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		rawToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if reflect.TypeOf(token.Method) == reflect.TypeOf(t.signingMethod) {
				return t.secretKey, nil
			}
			return nil, &errors.TokenSigningMethodMismatch{Message: token.Method.Alg()}
		},
	)
	if err != nil {
		return nil, &errors.ParseToken{Message: err.Error()}
	}
	userClaims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, &errors.InvalidClaims{Message: "could not load UserClaims from token"}
	}
	return userClaims, nil
}
