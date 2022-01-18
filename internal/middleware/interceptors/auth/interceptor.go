package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errs"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
)

type Interceptor struct {
	protectedRPCs map[string]string
	tokenManager  *tokenmanager.TokenManager
}

func New(tm *tokenmanager.TokenManager) *Interceptor {
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

func (i *Interceptor) check(ctx context.Context) (*tokenmanager.UserClaims, error) {
	metaData, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, &errs.Metadata{Message: "failed to load metadata from incoming context"}
	}
	authValues := metaData["authorization"]
	if len(authValues) < 1 {
		return nil, &errs.AuthHeader{Message: "unable to parse auth header from incoming context"}
	}
	// Do something with claims here eventually
	claims, err := i.tokenManager.Validate(authValues[0])
	if err != nil {
		return nil, &errs.TokenValidation{Message: err.Error()}
	}
	return claims, nil
}
