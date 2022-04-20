package tokenmanager

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	"reflect"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type Config struct {
	SecretKey     string
	SigningMethod string
	ValidDuration string
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
		return nil, &UnsupportedSigningMethod{Message: config.SigningMethod}
	}
	duration, err := time.ParseDuration(config.ValidDuration)
	if err != nil {
		return nil, &ParseDuration{Message: err.Error()}
	}
	return &TokenManager{
		secretKey:     config.SecretKey,
		signingMethod: method,
		validDuration: duration,
	}, nil
}

func (t *TokenManager) Generate(user *types.User) (string, error) {
	claims := &types.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.Email,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(t.validDuration)),
			ID:        strconv.Itoa(user.ID),
		},
	}
	token := jwt.NewWithClaims(t.signingMethod, claims)
	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", &TokenGeneration{Message: err.Error()}
	}
	return tokenString, nil
}

func (t *TokenManager) Validate(rawToken string) (*types.UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		rawToken,
		&types.UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if reflect.TypeOf(token.Method) == reflect.TypeOf(t.signingMethod) {
				return []byte(t.secretKey), nil
			}
			return nil, &TokenSigningMethodMismatch{Message: token.Method.Alg()}
		},
	)
	if err != nil {
		return nil, &ParseToken{Message: err.Error()}
	}
	userClaims, ok := token.Claims.(*types.UserClaims)
	if !ok {
		return nil, &InvalidClaims{Message: "could not load UserClaims from token"}
	}
	return userClaims, nil
}
