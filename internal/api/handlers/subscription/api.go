package subscription

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	subscriptionV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/subscription/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct{}

type Session interface {
	GetUserFromContext(context.Context) (*types.User, error)
}

type Server interface {
	CreateSubscription(int, int) error
	DeleteSubscription(int, int) error
}

type Handler struct {
	subscriptionV1.UnimplementedSubscriptionServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	session       Session
	server        Server
}

func New(c *Config, session Session, server Server, l *zap.Logger) *Handler {
	h := &Handler{
		log:           l,
		protectedRPCs: make(map[string]string),
		session:       session,
		server:        server,
	}
	h.requireAuth("Follow")
	h.requireAuth("Unfollow")
	return h
}

func (h *Handler) Follow(ctx context.Context, request *subscriptionV1.FollowRequest) (*subscriptionV1.FollowResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, Internal{}.Error())
	}
	err = h.server.CreateSubscription(user.ID, int(request.GetId()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, Internal{}.Error())
	case nil:
		return &subscriptionV1.FollowResponse{}, nil
	}
}

func (h *Handler) Unfollow(ctx context.Context, request *subscriptionV1.UnfollowRequest) (*subscriptionV1.UnfollowResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, Internal{}.Error())
	}
	err = h.server.DeleteSubscription(user.ID, int(request.GetId()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, Internal{}.Error())
	case nil:
		return &subscriptionV1.UnfollowResponse{}, nil
	}
}

func (h *Handler) Close() {}

func (h *Handler) GetProtectedRPCs() []string {
	protected := make([]string, len(h.protectedRPCs))
	for _, v := range h.protectedRPCs {
		protected = append(protected, v)
	}
	return protected
}

func (h *Handler) RegisterAPIServer(s *grpc.Server) {
	subscriptionV1.RegisterSubscriptionServiceServer(s, h)
}

func (h *Handler) requireAuth(rpcName string) {
	h.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		subscriptionV1.SubscriptionService_ServiceDesc.ServiceName,
		rpcName,
	)
}
