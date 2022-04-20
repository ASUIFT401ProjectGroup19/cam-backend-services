package comment

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	commentV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/comment/v1"
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
	Create(*types.User, *types.Comment) (*types.Comment, error)
	Read(int) (*types.Comment, error)
	Delete(*types.User, int) error
	Update(*types.User, *types.Comment) (*types.Comment, error)
	ReadByPostID(int) ([]*types.Comment, error)
}

type Handler struct {
	commentV1.UnimplementedCommentServiceServer
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
	h.requireAuth("Create")
	h.requireAuth("Update")
	h.requireAuth("Delete")
	return h
}

func (h *Handler) Create(ctx context.Context, request *commentV1.CreateRequest) (*commentV1.CreateResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "placeholder")
	}
	comment, err := h.server.Create(user, &types.Comment{
		Content:  request.GetComment().GetContent(),
		ParentID: int(request.GetComment().GetParentId()),
		PostID:   int(request.GetComment().GetPostId()),
		UserID:   user.ID,
	})
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		return &commentV1.CreateResponse{
			Id: int32(comment.ID),
		}, nil
	}
}

func (h *Handler) Read(ctx context.Context, request *commentV1.ReadRequest) (*commentV1.ReadResponse, error) {
	comment, err := h.server.Read(int(request.GetId()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		return &commentV1.ReadResponse{
			Comment: &commentV1.Comment{
				Id:       int32(comment.ID),
				Content:  comment.Content,
				ParentId: int32(comment.ParentID),
				PostId:   int32(comment.PostID),
				UserId:   int32(comment.UserID),
				UserName: comment.UserName,
			},
		}, nil
	}
}

func (h *Handler) Update(ctx context.Context, request *commentV1.UpdateRequest) (*commentV1.UpdateResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "placeholder")
	}
	comment, err := h.server.Update(user, &types.Comment{
		Content:  request.GetComment().GetContent(),
		ParentID: int(request.GetComment().GetParentId()),
		PostID:   int(request.GetComment().GetPostId()),
		UserID:   user.ID,
	})
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		return &commentV1.UpdateResponse{
			Id:      int32(comment.ID),
			Success: true,
		}, nil
	}
}

func (h *Handler) Delete(ctx context.Context, request *commentV1.DeleteRequest) (*commentV1.DeleteResponse, error) {
	user, err := h.session.GetUserFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "placeholder")
	}
	err = h.server.Delete(user, int(request.GetId()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		return &commentV1.DeleteResponse{
			Success: true,
		}, nil
	}
}

func (h *Handler) CommentsByPost(ctx context.Context, request *commentV1.CommentsByPostRequest) (*commentV1.CommentsByPostResponse, error) {
	comments, err := h.server.ReadByPostID(int(request.GetId()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		c := make([]*commentV1.Comment, len(comments))
		for i, v := range comments {
			c[i] = &commentV1.Comment{
				Id:       int32(v.ID),
				Content:  v.Content,
				ParentId: int32(v.ParentID),
				PostId:   int32(v.PostID),
				UserId:   int32(v.UserID),
				UserName: v.UserName,
			}
		}
		return &commentV1.CommentsByPostResponse{
			Comments: c,
		}, nil
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
	commentV1.RegisterCommentServiceServer(s, h)
}

func (h *Handler) requireAuth(rpcName string) {
	h.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		commentV1.CommentService_ServiceDesc.ServiceName,
		rpcName,
	)
}
