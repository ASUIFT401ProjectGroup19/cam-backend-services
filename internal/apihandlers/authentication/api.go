package authentication

import (
	"context"
	"fmt"

	authenticationAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/authentication/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errors"
	tm "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Storage interface {
	CheckPassword(username, password string) (*models.User, error)
	CreateUser(*models.User) (*models.User, error)
	RetrieveUserByID(int) (*models.User, error)
	RetrieveUserByUserName(string) (*models.User, error)
}

type Config struct{}

type APIv1 struct {
	authenticationAPIv1.UnimplementedAuthenticationServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	session       *tm.TokenManager
	storage       Storage
}

func New(config *Config, s Storage, log *zap.Logger, tm *tm.TokenManager) *APIv1 {
	return &APIv1{
		log:           log,
		protectedRPCs: make(map[string]string),
		storage:       s,
		session:       tm,
	}
}

func (a *APIv1) CreateAccount(
	ctx context.Context,
	request *authenticationAPIv1.CreateAccountRequest,
) (*authenticationAPIv1.CreateAccountResponse, error) {
	_, err := a.storage.CreateUser(
		&models.User{
			FirstName: request.GetFirstName(),
			LastName:  request.GetLastName(),
			Email:     request.GetUserName(),
			Password:  request.GetPassword(),
		},
	)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Unknown, err.Error())
	case *errors.BeginTransaction:
		return nil, status.Error(codes.Internal, err.Error())
	case *errors.EncryptPassword:
		return nil, status.Error(codes.Internal, err.Error())
	case *errors.Exists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case *errors.InsertRecord:
		return nil, status.Error(codes.Internal, err.Error())
	case nil:
		return &authenticationAPIv1.CreateAccountResponse{
			Success: true,
		}, nil
	}
}

func (a *APIv1) Login(
	ctx context.Context,
	request *authenticationAPIv1.LoginRequest,
) (*authenticationAPIv1.LoginResponse, error) {
	user, err := a.storage.CheckPassword(request.GetUserName(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	token, err := a.session.Generate(user)
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

func (a *APIv1) GetProtectedRPCs() []string {
	protected := make([]string, len(a.protectedRPCs))
	for _, v := range a.protectedRPCs {
		protected = append(protected, v)
	}
	return protected
}

func (a *APIv1) requireAuth(rpcName string) {
	a.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		authenticationAPIv1.AuthenticationService_ServiceDesc.ServiceName,
		rpcName,
	)
}
