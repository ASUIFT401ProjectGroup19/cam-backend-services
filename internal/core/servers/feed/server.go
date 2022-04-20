package feed

import (
	"errors"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	RetrieveMediaByPostID(int) ([]*types.Media, error)
	RetrieveSubscribedPostsPaginated(int, int, int) ([]*types.Post, error)
	ReadCommentsByPostID(int) ([]*types.Comment, error)
}

type Server struct {
	storage Storage
}

func New(storage Storage) *Server {
	s := &Server{
		storage: storage,
	}
	return s
}

func (s *Server) GetPostBatch(user *types.User, pageNumber int, batchSize int) ([]*types.Post, error) {
	posts, err := s.storage.RetrieveSubscribedPostsPaginated(user.ID, pageNumber, batchSize)
	if err != nil {
		return nil, errors.New("placeholder")
	}
	for i, v := range posts {
		media, err := s.storage.RetrieveMediaByPostID(v.ID)
		if err != nil {
			return nil, err
		}
		posts[i].Media = media
		comments, err := s.storage.ReadCommentsByPostID(v.ID)
		if err != nil {
			return nil, err
		}
		posts[i].Comments = comments
	}
	return posts, nil
}
