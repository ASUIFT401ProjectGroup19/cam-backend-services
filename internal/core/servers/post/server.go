package post

import (
	"context"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Session interface {
	GetUsernameFromContext(context.Context) (string, error)
}

type Storage interface {
	CreateMedia(*types.Media) (*types.Media, error)
	CreatePost(*types.Post) (*types.Post, error)
	RetrievePostByID(int) (*types.Post, error)
	RetrieveUserByUserName(string) (*types.User, error)
	RetrieveMediaByPostID(int) ([]*types.Media, error)
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

func (s *Server) Create(ctx context.Context, post *types.Post, media []*types.Media) (*types.Post, error) {
	username, _ := s.session.GetUsernameFromContext(ctx)
	user, _ := s.storage.RetrieveUserByUserName(username)
	post.UserID = user.ID
	postResult, _ := s.storage.CreatePost(post)
	for k := range media {
		media[k].PostID = postResult.ID
		_, _ = s.storage.CreateMedia(media[k])
	}
	return postResult, nil
}

func (s *Server) Read(id int) (*types.Post, []*types.Media, error) {
	post, err := s.storage.RetrievePostByID(id)
	if err != nil {
		return nil, nil, err
	}
	media, err := s.storage.RetrieveMediaByPostID(id)
	if err != nil {
		return nil, nil, err
	}
	return post, media, nil
}
