package identity

import (
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/middleware/tokenmanager"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Session interface {
	Generate(*models.User) (string, error)
	Validate(rawToken string) (*tokenmanager.UserClaims, error)
}

type Storage interface {
	CheckPassword(username, password string) (*models.User, error)
	CreateUser(*models.User) (*models.User, error)
	RetrieveUserByID(int) (*models.User, error)
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

func (s *Server) CreateAccount(user *models.User) (int, error) {
	createdUser, err := s.storage.CreateUser(user)
	if err != nil {
		return 0, err
	}
	return createdUser.ID, nil
}

func (s *Server) Login(username, password string) (string, error) {
	user, err := s.storage.CheckPassword(username, password)
	if err != nil {
		return "", status.Error(codes.PermissionDenied, err.Error())
	}
	token, err := s.session.Generate(user)
	if err != nil {
		return "", status.Error(codes.Internal, err.Error())
	}
	return token, nil
}
