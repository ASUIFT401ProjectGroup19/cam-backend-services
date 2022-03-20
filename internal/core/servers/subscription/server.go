package subscription

import (
	"context"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Session interface {
	GetUsernameFromContext(context.Context) (string, error)
}

type Storage interface {
	CreateSubscription(*types.User, *types.User) error
	DeleteSubscription(*types.User, *types.User) error
	RetrieveUserByID(int) (*types.User, error)
	RetrieveUserByUserName(string) (*types.User, error)
}

type Server struct {
	session Session
	storage Storage
}

func New(session Session, storage Storage) *Server {
	s := &Server{
		session: session,
		storage: storage,
	}
	return s
}

func (s *Server) CreateSubscription(ctx context.Context, id int) error {
	username, err := s.session.GetUsernameFromContext(ctx)
	if err != nil {
		return err
	}
	user, err := s.storage.RetrieveUserByUserName(username)
	if err != nil {
		return err
	}
	other, err := s.storage.RetrieveUserByID(id)
	if err != nil {
		return err
	}
	return s.storage.CreateSubscription(user, other)
}

func (s *Server) DeleteSubscription(ctx context.Context, id int) error {
	username, err := s.session.GetUsernameFromContext(ctx)
	if err != nil {
		return err
	}
	user, err := s.storage.RetrieveUserByUserName(username)
	if err != nil {
		return err
	}
	other, err := s.storage.RetrieveUserByID(id)
	if err != nil {
		return err
	}
	return s.storage.DeleteSubscription(user, other)
}
