package gallery

import (
	"errors"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	RetrieveMediaByPostID(int) ([]*types.Media, error)
	RetrieveUserPostsPaginated(int, int, int) ([]*types.Post, error)
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

func (s *Server) GetGallery(userID int, page int, batchSize int) ([]*types.Post, error) {
	posts, err := s.storage.RetrieveUserPostsPaginated(userID, page, batchSize)
	if err != nil {
		return nil, errors.New("placeholder")
	}
	for i, v := range posts {
		media, err := s.storage.RetrieveMediaByPostID(v.ID)
		if err != nil {
			return nil, err
		}
		posts[i].Media = media
	}
	return posts, nil
}