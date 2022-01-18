package identity

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/servers/identity"

	identityAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/identity/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errs"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Config struct{}

type APIv1 struct {
	identityAPIv1.UnimplementedIdentityServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	server        *identity.Server
}

func New(config *Config, s *identity.Server, log *zap.Logger) *APIv1 {
	return &APIv1{
		log:           log,
		protectedRPCs: make(map[string]string),
		server:        s,
	}
}

func (a *APIv1) CreateAccount(
	ctx context.Context,
	request *identityAPIv1.CreateAccountRequest,
) (*identityAPIv1.CreateAccountResponse, error) {
	_, err := a.server.CreateAccount(
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
	case *errs.BeginTransaction:
		return nil, status.Error(codes.Internal, err.Error())
	case *errs.EncryptPassword:
		return nil, status.Error(codes.Internal, err.Error())
	case *errs.Exists:
		return nil, status.Error(codes.AlreadyExists, err.Error())
	case *errs.InsertRecord:
		return nil, status.Error(codes.Internal, err.Error())
	case nil:
		return &identityAPIv1.CreateAccountResponse{
			Success: true,
		}, nil
	}
}

func (a *APIv1) Login(
	ctx context.Context,
	request *identityAPIv1.LoginRequest,
) (*identityAPIv1.LoginResponse, error) {
	token, err := a.server.Login(request.GetUserName(), request.GetPassword())
	switch err.(type) {
	default:
		return nil, status.Error(codes.Unknown, err.Error())
	case nil:
		return &identityAPIv1.LoginResponse{
			Token: token,
		}, nil
	}
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	identityAPIv1.RegisterIdentityServiceServer(server, a)
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
		identityAPIv1.IdentityService_ServiceDesc.ServiceName,
		rpcName,
	)
}
