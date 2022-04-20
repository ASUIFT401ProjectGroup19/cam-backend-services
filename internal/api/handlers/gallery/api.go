package gallery

import (
	"context"
	"fmt"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	commentV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/comment/v1"
	galleryV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/gallery/v1"
	postV1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Config struct{}

type Session interface {
	GetUserFromContext(ctx context.Context) (*types.User, error)
}

type Server interface {
	GetGallery(int, int, int) ([]*types.Post, error)
}

type Handler struct {
	galleryV1.UnimplementedGalleryServiceServer
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
	return h
}

func (h *Handler) Gallery(ctx context.Context, request *galleryV1.GalleryRequest) (*galleryV1.GalleryResponse, error) {
	posts, err := h.server.GetGallery(int(request.UserId), int(request.GetPage()), int(request.GetBatchSize()))
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, "placeholder")
	case nil:
		response := &galleryV1.GalleryResponse{}
		response.Posts = make([]*postV1.Post, len(posts))
		for i, v := range posts {
			response.Posts[i] = &postV1.Post{
				Id:          int32(v.ID),
				Description: v.Description,
				UserId:      int32(v.UserID),
				UserName:    v.UserName,
			}
			response.Posts[i].Media = make([]*postV1.Media, len(v.Media))
			for ii, vv := range v.Media {
				response.Posts[i].Media[ii] = &postV1.Media{
					Link: vv.Link,
				}
			}
			response.Posts[i].Comments = make([]*commentV1.Comment, len(v.Comments))
			for ii, vv := range v.Comments {
				response.Posts[i].Comments[ii] = &commentV1.Comment{
					Id:       int32(vv.ID),
					Content:  vv.Content,
					ParentId: int32(vv.ParentID),
					PostId:   int32(vv.PostID),
					UserId:   int32(vv.UserID),
					UserName: vv.UserName,
				}
			}
		}
		return response, nil
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
	galleryV1.RegisterGalleryServiceServer(s, h)
}

func (h *Handler) requireAuth(rpcName string) {
	h.protectedRPCs[rpcName] = fmt.Sprintf(
		"/%s/%s",
		galleryV1.GalleryService_ServiceDesc.ServiceName,
		rpcName,
	)
}
