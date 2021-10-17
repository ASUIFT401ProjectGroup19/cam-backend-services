package v1

import (
	"context"
	authenticationAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-backend-services/pkg/gen/go/grpc/authentication/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {}

type APIv1 struct {
	authenticationAPIv1.UnimplementedAuthenticationAPIServer
	log *zap.Logger
}

func (a *APIv1) CreateAccount(ctx context.Context, request *authenticationAPIv1.CreateAccountRequest) (*authenticationAPIv1.CreateAccountResponse, error) {
	a.log.Info(request.AccountName)
	resp := &authenticationAPIv1.CreateAccountResponse{Success: true}
	return resp, nil
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	authenticationAPIv1.RegisterAuthenticationAPIServer(server, a)
}

func New(logger *zap.Logger) *APIv1 {
	return &APIv1{
		log: logger,
	}
}
