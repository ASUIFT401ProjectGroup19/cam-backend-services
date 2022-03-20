package feed

import (
	"context"
	"fmt"
	feedV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/feed/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct {
}

type Server interface {
}

type Handler struct {
	feedV1.UnimplementedFeedServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	server        Server
}

func New(c *Config, s Server, l *zap.Logger) *Handler {
	h := &Handler{
		log:           l,
		protectedRPCs: make(map[string]string),
		server:        s,
	}
	return h
}

func (h *Handler) Next(context.Context, *feedV1.NextRequest) (*feedV1.NextResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NYI")
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
	feedV1.RegisterFeedServiceServer(s, h)
}

func (h *Handler) requireAuth(rpcName string) {
	h.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		feedV1.FeedService_ServiceDesc.ServiceName,
		rpcName,
	)
}
