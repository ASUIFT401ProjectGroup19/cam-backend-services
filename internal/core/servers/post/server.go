package post

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	CreateMedia(*types.Media) (*types.Media, error)
	CreatePost(*types.Post) (*types.Post, error)
	RetrievePostByID(int) (*types.Post, error)
	RetrieveUserByUserName(string) (*types.User, error)
	RetrieveMediaByPostID(int) ([]*types.Media, error)
	ReadCommentsByPostID(int) ([]*types.Comment, error)
}

type Server struct {
	storage Storage
}

func New(storage Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) Create(user *types.User, post *types.Post, media []*types.Media, comments []*types.Comment) (*types.Post, error) {
	post.UserID = user.ID
	postResult, _ := s.storage.CreatePost(post)
	for k := range media {
		media[k].PostID = postResult.ID
		_, _ = s.storage.CreateMedia(media[k])
	}
	for k := range comments {
		comments[k].PostID = postResult.ID
		comments[k].UserID = postResult.UserID
	}
	return postResult, nil
}

func (s *Server) Read(id int) (*types.Post, error) {
	post, err := s.storage.RetrievePostByID(id)
	if err != nil {
		return nil, err
	}
	media, err := s.storage.RetrieveMediaByPostID(id)
	if err != nil {
		return nil, err
	}
	post.Media = media
	comments, err := s.storage.ReadCommentsByPostID(id)
	if err != nil {
		return nil, err
	}
	post.Comments = comments
	return post, nil
}
