package post

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/servers/post"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	commentV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/comment/v1"
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
	Create(*types.User, *types.Post, []*types.Media, []*types.Comment) (*types.Post, error)
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
	comments := make([]*types.Comment, len(req.GetPost().GetComments()))
	for k, v := range req.GetPost().GetComments() {
		comments[k] = &types.Comment{
			Content:  v.Content,
			ParentID: int(v.ParentId),
			PostID:   int(v.PostId),
			UserID:   user.ID,
		}
	}
	p, err := h.server.Create(user,
		&types.Post{
			Description: req.GetPost().GetDescription(),
		},
		media,
		comments,
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
	result, err := h.server.Read(int(req.GetId()))
	if err != nil {
		return nil, err
	}
	media := make([]*postV1.Media, len(result.Media))
	for k, v := range result.Media {
		media[k] = &postV1.Media{
			Link: v.Link,
		}
	}
	comments := make([]*commentV1.Comment, len(result.Comments))
	for i, v := range result.Comments {
		comments[i] = &commentV1.Comment{
			Id:       int32(v.ID),
			Content:  v.Content,
			ParentId: int32(v.ParentID),
			PostId:   int32(v.PostID),
			UserId:   int32(v.UserID),
			UserName: v.UserName,
		}
	}
	return &postV1.ReadResponse{
		Post: &postV1.Post{
			Id:          int32(result.ID),
			Description: result.Description,
			Media:       media,
			Comments:    comments,
			UserId:      int32(result.UserID),
			UserName:    result.UserName,
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
