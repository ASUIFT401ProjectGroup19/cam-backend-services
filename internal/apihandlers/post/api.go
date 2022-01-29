package post

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/servers/post"

	postAPIv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Config struct{}

type APIv1 struct {
	postAPIv1.UnimplementedPostServiceServer
	log           *zap.Logger
	server        *post.Server
	protectedRPCs map[string]string
}

func New(config *Config, s *post.Server, log *zap.Logger) *APIv1 {
	a := &APIv1{
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

func (a *APIv1) Create(ctx context.Context, req *postAPIv1.CreateRequest) (*postAPIv1.CreateResponse, error) {
	media := make([]*models.Media, len(req.GetPost().GetMedia()))
	for k, v := range req.GetPost().GetMedia() {
		media[k] = &models.Media{
			Link: v.GetLink(),
		}
	}
	p, err := a.server.Create(ctx,
		&models.Post{
			Description: req.GetPost().GetDescription(),
		},
		media,
	)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, err.Error())
	case nil:
		return &postAPIv1.CreateResponse{
			Id: int32(p.ID),
		}, nil
	}
}

func (a *APIv1) Read(ctx context.Context, req *postAPIv1.ReadRequest) (*postAPIv1.ReadResponse, error) {
	postResponse, mediaResponse, err := a.server.Read(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	media := make([]*postAPIv1.Media, len(mediaResponse))
	for k, v := range mediaResponse {
		media[k] = &postAPIv1.Media{
			Link: v.Link,
		}
	}
	return &postAPIv1.ReadResponse{
		Post: &postAPIv1.Post{
			Id:          int32(postResponse.ID),
			Description: postResponse.Description,
			Media:       media,
		},
	}, nil
}

func (a *APIv1) Update(context.Context, *postAPIv1.UpdateRequest) (*postAPIv1.UpdateResponse, error) {
	return &postAPIv1.UpdateResponse{}, nil
}

func (a *APIv1) Delete(context.Context, *postAPIv1.DeleteRequest) (*postAPIv1.DeleteResponse, error) {
	return &postAPIv1.DeleteResponse{}, nil
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	postAPIv1.RegisterPostServiceServer(server, a)
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
		postAPIv1.PostService_ServiceDesc.ServiceName,
		rpcName,
	)
}
