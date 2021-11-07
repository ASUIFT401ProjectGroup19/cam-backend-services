package authentication

import (
	"context"

	authenticationAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/authentication/v1"
	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	driver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/database/cam"
	tm "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
)

type Config struct{}

type APIv1 struct {
	authenticationAPIv1.UnimplementedAuthenticationServiceServer
	db  *driver.Driver
	log *zap.Logger
	tm  *tm.TokenManager
}

func New(config *Config, db *driver.Driver, log *zap.Logger, tm *tm.TokenManager) *APIv1 {
	return &APIv1{
		db:  db,
		log: log,
		tm:  tm,
	}
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
	case *driver.ErrorBeginTransaction:
		return nil, status.Error(codes.Internal, err.Error())
	case *driver.ErrorEncryptPassword:
		return nil, status.Error(codes.Internal, err.Error())
	case *driver.ErrorExists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case *driver.ErrorInsertRecord:
		return nil, status.Error(codes.Internal, err.Error())
	case *driver.ErrorUnknown:
		return nil, status.Error(codes.Unknown, err.Error())
	}
	return &authenticationAPIv1.CreateAccountResponse{
		Success: true,
	}, nil
}

func (a *APIv1) Login(ctx context.Context, request *authenticationAPIv1.LoginRequest) (*authenticationAPIv1.LoginResponse, error) {
	user, err := a.db.GetUser(request.GetUserName())
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	err = driver.CheckPassword(user, request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	token, err := a.tm.Generate(user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &authenticationAPIv1.LoginResponse{
		Token: token,
	}, nil
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	authenticationAPIv1.RegisterAuthenticationServiceServer(server, a)
}
