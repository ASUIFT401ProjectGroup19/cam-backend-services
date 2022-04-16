package post

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/post"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	postV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
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
	Create(*types.User, *types.Post, []*types.Media) (*types.Post, error)
	Read(int) (*types.Post, error)
}

type Handler struct {
	postV1.UnimplementedPostServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	session       Session
	server        Server
}

func New(config *Config, session Session, server *post.Server, log *zap.Logger) *Handler {
	h := &Handler{
		log:           log,
		session:       session,
		server:        server,
		protectedRPCs: make(map[string]string),
	}
	h.requireAuth("Create")
	h.requireAuth("Read")
	h.requireAuth("Update")
	h.requireAuth("Delete")
	return h
}

func (h *Handler) Create(ctx context.Context, req *postV1.CreateRequest) (*postV1.CreateResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	media := make([]*types.Media, len(req.GetPost().GetMedia()))
	for k, v := range req.GetPost().GetMedia() {
		media[k] = &types.Media{
			Link: v.GetLink(),
		}
	}
	p, err := h.server.Create(user,
		&types.Post{
			Description: req.GetPost().GetDescription(),
		},
		media,
	)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, err.Error())
	case nil:
		return &postV1.CreateResponse{
			Id: int32(p.ID),
		}, nil
	}
}

func (h *Handler) Read(ctx context.Context, req *postV1.ReadRequest) (*postV1.ReadResponse, error) {
	postResponse, err := h.server.Read(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	media := make([]*postV1.Media, len(postResponse.Media))
	for k, v := range postResponse.Media {
		media[k] = &postV1.Media{
			Link: v.Link,
		}
	}
	return &postV1.ReadResponse{
		Post: &postV1.Post{
			Id:          int32(postResponse.ID),
			Description: postResponse.Description,
			Media:       media,
		},
	}, nil
}

func (h *Handler) Update(context.Context, *postV1.UpdateRequest) (*postV1.UpdateResponse, error) {
	return &postV1.UpdateResponse{}, nil
}

func (h *Handler) Delete(context.Context, *postV1.DeleteRequest) (*postV1.DeleteResponse, error) {
	return &postV1.DeleteResponse{}, nil
}

func (h *Handler) Close() {}

func (h *Handler) RegisterAPIServer(server *grpc.Server) {
	postV1.RegisterPostServiceServer(server, h)
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
		postV1.PostService_ServiceDesc.ServiceName,
		rpcName,
	)
}
