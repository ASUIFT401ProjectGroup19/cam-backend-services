package post

import (
	"context"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

type Session interface {
	GetUsernameFromContext(context.Context) (string, error)
}

type Storage interface {
	CreateMedia(*models.Media) (*models.Media, error)
	CreatePost(*models.Post) (*models.Post, error)
	RetrievePostByID(int) (*models.Post, error)
	RetrieveUserByUserName(string) (*models.User, error)
}

type Server struct {
	session Session
	storage Storage
}

func New(session Session, storage Storage) *Server {
	return &Server{
		session: session,
		storage: storage,
	}
}

func (s *Server) Create(ctx context.Context, post *models.Post, media *models.Media) (*models.Post, *models.Media, error) {
	username, _ := s.session.GetUsernameFromContext(ctx)
	user, _ := s.storage.RetrieveUserByUserName(username)
	post.UserID = user.ID
	postResult, _ := s.storage.CreatePost(post)
	media.PostID = postResult.ID
	mediaResult, _ := s.storage.CreateMedia(media)
	return postResult, mediaResult, nil
}
