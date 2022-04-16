package session

import (
	"context"
	"errors"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/session/tokenmanager"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

type Storage interface {
	RetrieveUserByUserName(string) (*types.User, error)
}

type Manager struct {
	tokenmanager.TokenManager
	storage Storage
}

func New(storage Storage, tm *tokenmanager.TokenManager) *Manager {
	return &Manager{
		TokenManager: *tm,
		storage:      storage,
	}
}

func (m *Manager) GetUserFromContext(ctx context.Context) (*types.User, error) {
	if claims, ok := ctx.Value("claims").(*types.UserClaims); !ok {
		return nil, errors.New("placeholder")
	} else {
		user, err := m.storage.RetrieveUserByUserName(claims.Subject)
		if err != nil {
			return nil, errors.New("placeholder")
		}
		return user, nil
	}
}
