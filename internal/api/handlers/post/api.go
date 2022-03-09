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

type Server interface {
	Create(context.Context, *types.Post, []*types.Media) (*types.Post, error)
	Read(int) (*types.Post, []*types.Media, error)
}

type Handler struct {
	postV1.UnimplementedPostServiceServer
	log           *zap.Logger
	protectedRPCs map[string]string
	server        Server
}

func New(config *Config, s *post.Server, log *zap.Logger) *Handler {
	a := &Handler{
		log:           log,
		server:        s,
		protectedRPCs: make(map[string]string),
	}
	a.requireAuth("Create")
	a.requireAuth("Read")
	a.requireAuth("Update")
	a.requireAuth("Delete")
	return a
}

func (a *Handler) Create(ctx context.Context, req *postV1.CreateRequest) (*postV1.CreateResponse, error) {
	media := make([]*types.Media, len(req.GetPost().GetMedia()))
	for k, v := range req.GetPost().GetMedia() {
		media[k] = &types.Media{
			Link: v.GetLink(),
		}
	}
	p, err := a.server.Create(ctx,
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

func (a *Handler) Read(ctx context.Context, req *postV1.ReadRequest) (*postV1.ReadResponse, error) {
	postResponse, mediaResponse, err := a.server.Read(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	media := make([]*postV1.Media, len(mediaResponse))
	for k, v := range mediaResponse {
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

func (a *Handler) Update(context.Context, *postV1.UpdateRequest) (*postV1.UpdateResponse, error) {
	return &postV1.UpdateResponse{}, nil
}

func (a *Handler) Delete(context.Context, *postV1.DeleteRequest) (*postV1.DeleteResponse, error) {
	return &postV1.DeleteResponse{}, nil
}

func (a *Handler) Close() {}

func (a *Handler) RegisterAPIServer(server *grpc.Server) {
	postV1.RegisterPostServiceServer(server, a)
}

func (a *Handler) GetProtectedRPCs() []string {
	protected := make([]string, len(a.protectedRPCs))
	for _, v := range a.protectedRPCs {
		protected = append(protected, v)
	}
	return protected
}

func (a *Handler) requireAuth(rpcName string) {
	a.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		postV1.PostService_ServiceDesc.ServiceName,
		rpcName,
	)
}