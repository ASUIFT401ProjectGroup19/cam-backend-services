package identity

import (
	"context"
	"errors"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/persistence/cam/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/session/tokenmanager"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	identityV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/identity/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct{}

type Session interface {
	GetUserFromContext(context.Context) (*types.User, error)
	Generate(*types.User) (string, error)
}

type Server interface {
	CreateAccount(*types.User) (int, error)
	Login(string, string) (*types.User, error)
}

type Handler struct {
	identityV1.UnimplementedIdentityServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	session       Session
	server        Server
}

func New(config *Config, session Session, server Server, log *zap.Logger) *Handler {
	h := &Handler{
		log:           log,
		protectedRPCs: make(map[string]string),
		session:       session,
		server:        server,
	}
	h.requireAuth("Refresh")
	return h
}

func (h *Handler) CreateAccount(ctx context.Context, request *identityV1.CreateAccountRequest) (*identityV1.CreateAccountResponse, error) {
	_, err := h.server.CreateAccount(
		&types.User{
			FirstName: request.GetFirstName(),
			LastName:  request.GetLastName(),
			Email:     request.GetUserName(),
			Password:  request.GetPassword(),
		},
	)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Unknown, Unknown{}.Error())
	case *cam.Exists:
		return nil, status.Error(codes.AlreadyExists, AccountExists{}.Error())
	case nil:
		return &identityV1.CreateAccountResponse{
			Success: true,
		}, nil
	}
}

func (h *Handler) Login(ctx context.Context, request *identityV1.LoginRequest) (*identityV1.LoginResponse, error) {
	user, err := h.server.Login(request.GetUserName(), request.GetPassword())
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, LoginFailed{}.Error())
	}
	token, err := h.session.Generate(user)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Unknown, Unknown{}.Error())
	case *tokenmanager.TokenGeneration:
		return nil, status.Error(codes.Internal, Internal{}.Error())
	case nil:
		return &identityV1.LoginResponse{
			Token: token,
		}, nil
	}
}

func (h *Handler) Refresh(ctx context.Context, request *identityV1.RefreshRequest) (*identityV1.RefreshResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, errors.New("placeholder")
	}
	token, err := h.session.Generate(user)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Unknown, Unknown{}.Error())
	case *tokenmanager.TokenGeneration:
		return nil, status.Error(codes.Internal, Internal{}.Error())
	case nil:
		return &identityV1.RefreshResponse{
			Token: token,
		}, nil
	}
}

func (h *Handler) Close() {}

func (h *Handler) RegisterAPIServer(server *grpc.Server) {
	identityV1.RegisterIdentityServiceServer(server, h)
}

func (h *Handler) GetProtectedRPCs() []string {
	protected := make([]string, len(h.protectedRPCs))
	for _, v := range h.protectedRPCs {
		protected = append(protected, v)
	}
	return protected
}

func (h *Handler) requireAuth(rpcName string) {
	h.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		identityV1.IdentityService_ServiceDesc.ServiceName,
		rpcName,
	)
}
