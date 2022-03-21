package auth

import (
	"context"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type TokenManager interface {
	Validate(string) (*types.UserClaims, error)
}

type Interceptor struct {
	protectedRPCs map[string]string
	tokenManager  TokenManager
}

func New(tm TokenManager) *Interceptor {
	return &Interceptor{
		protectedRPCs: make(map[string]string),
		tokenManager:  tm,
	}
}

func (i *Interceptor) RegisterProtectedRoutes(routes []string) {
	for _, e := range routes {
		i.protectedRPCs[e] = ""
	}
}

func (i *Interceptor) Stream() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		if _, ok := i.protectedRPCs[info.FullMethod]; ok {
			if _, err := i.check(ss.Context()); err != nil {
				return err
			}
		}
		return handler(srv, ss)
	}
}

func (i *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if _, ok := i.protectedRPCs[info.FullMethod]; ok {
			claims, err := i.check(ctx)
			if err != nil {
				return nil, err
			}
			return handler(context.WithValue(ctx, "claims", claims), req)
		}
		return handler(ctx, req)
	}
}

func (i *Interceptor) check(ctx context.Context) (*types.UserClaims, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, &Metadata{Message: "failed to load metadata from incoming context"}
	}
	authValues := metaData["authorization"]
	if len(authValues) < 1 {
		return nil, &MissingHeader{Message: "unable to parse auth header from incoming context"}
	}
	token := authValues[0]
	if token[:7] == "Bearer " {
		token = token[7:]
	}
	claims, err := i.tokenManager.Validate(token)
	if err != nil {
		return nil, &TokenValidation{Message: err.Error()}
	}
	return claims, nil
}
