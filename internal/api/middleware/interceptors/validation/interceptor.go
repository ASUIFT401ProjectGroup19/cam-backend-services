package validation

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validator interface {
	ValidateAll() error
}

type Interceptor struct{}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if err := validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func (i *Interceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// We don't have any streaming RPCs yet so this is plumbing for later.
		return handler(srv, ss)
	}
}

func validate(req interface{}) error {
	switch request := req.(type) {
	case validator:
		if err := request.ValidateAll(); err != nil {
			return status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return nil
}
