package comment

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	CreateComment(*types.Comment) (*types.Comment, error)
	ReadComment(int) (*types.Comment, error)
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

func (s *Server) Create(user *types.User, comment *types.Comment) (*types.Comment, error) {
	comment.UserID = user.ID
	result, err := s.storage.CreateComment(comment)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *Server) Read(commentID int) (*types.Comment, error) {
	comment, err := s.storage.ReadComment(commentID)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *Server) Delete(user *types.User, commentID int) error {
	return nil
}

func (s *Server) Update(user *types.User, commentID *types.Comment) (*types.Comment, error) {
	return nil, nil
}

func (s *Server) ReadByPostID(postID int) ([]*types.Comment, error) {
	comments, err := s.storage.ReadCommentsByPostID(postID)
	if err != nil {
		return nil, err
	}
	return comments, nil
}
