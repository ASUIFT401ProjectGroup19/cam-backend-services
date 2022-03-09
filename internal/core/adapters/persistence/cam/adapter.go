package cam

import (
	camDriver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/persistence/cam/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

// Adapter implements identity.Storage.
// Adapter implements post.Storage.
type Adapter struct {
	driver *camDriver.Driver
}

func New(driver *camDriver.Driver) *Adapter {
	return &Adapter{
		driver: driver,
	}
}

func (a *Adapter) CheckPassword(username string, password string) (*types.User, error) {
	user, err := a.driver.RetrieveUserByUserName(username)
	if err != nil {
		return nil, err
	}
	err = a.driver.CheckPassword(user, password)
	if err != nil {
		return nil, err
	}
	return userDriverToModel(user), nil
}

func (a *Adapter) CreateUser(user *types.User) (*types.User, error) {
	u, err := a.driver.CreateUser(userModelToDriver(user))
	if err != nil {
		return nil, err
	}
	return userDriverToModel(u), nil
}

func (a *Adapter) RetrieveUserByID(id int) (*types.User, error) {
	user, err := a.driver.RetrieveUserByID(id)
	if err != nil {
		return nil, err
	}
	return userDriverToModel(user), nil
}

func (a *Adapter) RetrieveUserByUserName(username string) (*types.User, error) {
	user, err := a.driver.RetrieveUserByUserName(username)
	if err != nil {
		return nil, err
	}
	return userDriverToModel(user), nil
}

func (a *Adapter) CreateMedia(media *types.Media) (*types.Media, error) {
	m, err := a.driver.CreateMedia(mediaModelToDriver(media))
	if err != nil {
		return nil, err
	}
	return mediaDriverToModel(m), nil
}

func (a *Adapter) CreatePost(post *types.Post) (*types.Post, error) {
	p, err := a.driver.CreatePost(postModelToDriver(post))
	if err != nil {
		return nil, err
	}
	return postDriverToModel(p), nil
}

func (a *Adapter) RetrievePostByID(i int) (*types.Post, error) {
	p, err := a.driver.RetrievePostByID(i)
	if err != nil {
		return nil, err
	}
	return postDriverToModel(p), nil
}

func (a *Adapter) RetrieveMediaByPostID(id int) ([]*types.Media, error) {
	m, err := a.driver.RetrieveMediaByPostID(id)
	if err != nil {
		return nil, err
	}
	media := make([]*types.Media, len(m))
	for k, v := range m {
		media[k] = mediaDriverToModel(v)
	}
	return media, err
}
