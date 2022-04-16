package subscription

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	CreateSubscription(*types.User, *types.User) error
	DeleteSubscription(*types.User, *types.User) error
	RetrieveUserByID(int) (*types.User, error)
	RetrieveUserByUserName(string) (*types.User, error)
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

func (s *Server) CreateSubscription(userID, otherID int) error {
	user, err := s.storage.RetrieveUserByID(userID)
	if err != nil {
		return err
	}
	other, err := s.storage.RetrieveUserByID(otherID)
	if err != nil {
		return err
	}
	return s.storage.CreateSubscription(user, other)
}

func (s *Server) DeleteSubscription(userID, otherID int) error {
	user, err := s.storage.RetrieveUserByID(userID)
	if err != nil {
		return err
	}
	other, err := s.storage.RetrieveUserByID(otherID)
	if err != nil {
		return err
	}
	return s.storage.DeleteSubscription(user, other)
}
