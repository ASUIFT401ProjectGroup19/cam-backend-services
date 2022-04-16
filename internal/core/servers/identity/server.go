package identity

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	CheckPassword(username, password string) (*types.User, error)
	CreateUser(*types.User) (*types.User, error)
	RetrieveUserByUserName(string) (*types.User, error)
}

type Server struct {
	storage Storage
}

func New(storage Storage) *Server {
	return &Server{
		storage: storage,
	}
}

func (s *Server) CreateAccount(user *types.User) (int, error) {
	createdUser, err := s.storage.CreateUser(user)
	if err != nil {
		return 0, err
	}
	return createdUser.ID, nil
}

func (s *Server) Login(username, password string) (*types.User, error) {
	user, err := s.storage.CheckPassword(username, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}
