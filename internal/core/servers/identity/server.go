package identity

import (
	"context"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Session interface {
	GetUsernameFromContext(context.Context) (string, error)
	Generate(*types.User) (string, error)
}

type Storage interface {
	CheckPassword(username, password string) (*types.User, error)
	CreateUser(*types.User) (*types.User, error)
	RetrieveUserByUserName(string) (*types.User, error)
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

func (s *Server) GenerateToken(user *types.User) (string, error) {
	token, err := s.session.Generate(user)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Server) RefreshToken(ctx context.Context) (string, error) {
	username, err := s.session.GetUsernameFromContext(ctx)
	if err != nil {
		return "", err
	}
	user, err := s.storage.RetrieveUserByUserName(username)
	if err != nil {
		return "", err
	}
	token, err := s.session.Generate(user)
	if err != nil {
		return "", err
	}
	return token, nil
}
