package authentication

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/camdb"

	authenticationAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/authentication/v1"
	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct{}

type APIv1 struct {
	authenticationAPIv1.UnimplementedAuthenticationServiceServer
	db  *camdb.DB
	log *zap.Logger
}

func (a *APIv1) CreateAccount(ctx context.Context, request *authenticationAPIv1.CreateAccountRequest) (*authenticationAPIv1.CreateAccountResponse, error) {
	if err := request.ValidateAll(); err != nil {
		a.log.Error("validating request", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user := &cam.User{
		FirstName: request.GetFirstName(),
		LastName:  request.GetLastName(),
		Email:     request.GetUserName(),
		Password:  request.GetPassword(),
	}
	err := a.db.SetUser(user)
	switch err.(type) {
	case *camdb.ErrorBeginTransaction:
		return nil, status.Error(codes.Internal, err.Error())
	case *camdb.ErrorEncryptPassword:
		return nil, status.Error(codes.Internal, err.Error())
	case *camdb.ErrorExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case *camdb.ErrorInsertRecord:
		return nil, status.Error(codes.Internal, err.Error())
	case *camdb.ErrorUnknown:
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &authenticationAPIv1.CreateAccountResponse{
		Success: true,
	}, nil
}

func (a *APIv1) Login(ctx context.Context, request *authenticationAPIv1.LoginRequest) (*authenticationAPIv1.LoginResponse, error) {
	return nil, nil
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	authenticationAPIv1.RegisterAuthenticationServiceServer(server, a)
}

func New(config *Config, db *camdb.DB, logger *zap.Logger) *APIv1 {
	return &APIv1{
		db:  db,
		log: logger,
	}
}
