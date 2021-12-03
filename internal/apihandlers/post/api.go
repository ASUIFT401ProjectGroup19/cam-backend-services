package post

import (
	"context"
	"fmt"

	postv1 "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/proto/go/post/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Storage interface {
	CreateMedia(*models.Media) (*models.Media, error)
	CreatePost(*models.Post) (*models.Post, error)
	RetrievePostByID(int) (*models.Post, error)
	RetrieveUserByUserName(string) (*models.User, error)
}

type Config struct{}

type APIv1 struct {
	postv1.UnimplementedPostServiceServer
	log           *zap.Logger
	storage       Storage
	protectedRPCs map[string]string
}

func New(config *Config, s Storage, log *zap.Logger) *APIv1 {
	a := &APIv1{
		log:           log,
		storage:       s,
		protectedRPCs: make(map[string]string),
	}
	a.requireAuth("Create")
	a.requireAuth("Read")
	a.requireAuth("Update")
	a.requireAuth("Delete")
	return a
}

func (a *APIv1) Create(ctx context.Context, req *postv1.CreateRequest) (*postv1.CreateResponse, error) {
	claims := ctx.Value("claims").(*tokenmanager.UserClaims)
	user, err := a.storage.RetrieveUserByUserName(claims.Subject)
	postReq := req.GetPost()
	mediaReq := req.GetPost().GetMedia()
	post, err := a.storage.CreatePost(
		&models.Post{
			Description: postReq.GetDescription(),
			UserID:      user.ID,
		},
	)
	_, _ = a.storage.CreateMedia(
		&models.Media{
			Link:   mediaReq.GetLink(),
			PostID: post.ID,
		},
	)
	switch err.(type) {
	default:
		return nil, status.Error(codes.Internal, err.Error())
	case nil:
		return &postv1.CreateResponse{
			Id: int32(post.ID),
		}, nil
	}
}

func (a *APIv1) Read(context.Context, *postv1.ReadRequest) (*postv1.ReadResponse, error) {
	return nil, nil
}

func (a *APIv1) Update(context.Context, *postv1.UpdateRequest) (*postv1.UpdateResponse, error) {
	return nil, nil
}

func (a *APIv1) Delete(context.Context, *postv1.DeleteRequest) (*postv1.DeleteResponse, error) {
	return nil, nil
}

func (a *APIv1) Close() {}

func (a *APIv1) RegisterAPIServer(server *grpc.Server) {
	postv1.RegisterPostServiceServer(server, a)
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
		postv1.PostService_ServiceDesc.ServiceName,
		rpcName,
	)
}
