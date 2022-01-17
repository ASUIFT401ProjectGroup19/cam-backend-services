package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
)

type AuthInterceptor struct {
	protectedRoutes map[string]string
	tokenManager    *tokenmanager.TokenManager
}

func New(tm *tokenmanager.TokenManager) *AuthInterceptor {
	return &AuthInterceptor{
		tokenManager: tm,
	}
}

func (a *AuthInterceptor) RegisterProtectedRoutes(routes []string) {
	for _, e := range routes {
		a.protectedRoutes[e] = ""
	}
}

func (a *AuthInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if _, ok := a.protectedRoutes[info.FullMethod]; ok {
			if err := a.check(ss.Context()); err != nil {
				return err
			}
		}
		return handler(srv, ss)
	}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if _, ok := a.protectedRoutes[info.FullMethod]; ok {
			if err := a.check(ctx); err != nil {
				return nil, err
			}
		}
		return handler(ctx, req)
	}
}

func (a *AuthInterceptor) check(ctx context.Context) error {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &ErrorMetadata{msg: "failed to load metadata from incoming context"}
	}
	authValues := metaData["authorization"]
	if len(authValues) < 1 {
		return &ErrorAuthHeader{msg: "unable to parse auth header from incoming context"}
	}
	// Do something with claims here eventually
	_, err := a.tokenManager.Validate(authValues[0])
	if err != nil {
		return &ErrorTokenValidation{msg: err.Error()}
	}
	return nil
}
