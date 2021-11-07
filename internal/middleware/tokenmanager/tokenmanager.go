package tokenmanager

import (
	"reflect"
	"time"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"

	"github.com/golang-jwt/jwt/v4"
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
		return nil, &ErrorUnsupportedSigningMethod{msg: config.SigningMethod}
	}
	duration, err := time.ParseDuration(config.ValidDuration)
	if err != nil {
		return nil, &ErrorParseDuration{msg: err.Error()}
	}
	return &TokenManager{
		secretKey:     config.SecretKey,
		signingMethod: method,
		validDuration: duration,
	}, nil
}

func (t *TokenManager) Generate(user *cam.User) (string, error) {
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
			return nil, &ErrorTokenSigningMethodMismatch{msg: token.Method.Alg()}
		},
	)
	if err != nil {
		return nil, &ErrorParseToken{msg: err.Error()}
	}
	userClaims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, &ErrorInvalidClaims{msg: "could not load UserClaims from token"}
	}
	return userClaims, nil
}
