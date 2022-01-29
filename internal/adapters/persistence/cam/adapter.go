package cam

import (
	camDriver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
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

func (a *Adapter) CheckPassword(username string, password string) (*models.User, error) {
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

func (a *Adapter) CreateUser(user *models.User) (*models.User, error) {
	u, err := a.driver.CreateUser(userModelToDriver(user))
	if err != nil {
		return nil, err
	}
	return userDriverToModel(u), nil
}

func (a *Adapter) RetrieveUserByID(i int) (*models.User, error) {
	// TODO implement me
	panic("implement me")
}

func (a *Adapter) RetrieveUserByUserName(username string) (*models.User, error) {
	user, err := a.driver.RetrieveUserByUserName(username)
	if err != nil {
		return nil, err
	}
	return userDriverToModel(user), nil
}

func (a *Adapter) CreateMedia(media *models.Media) (*models.Media, error) {
	m, err := a.driver.CreateMedia(mediaModelToDriver(media))
	if err != nil {
		return nil, err
	}
	return mediaDriverToModel(m), nil
}

func (a *Adapter) CreatePost(post *models.Post) (*models.Post, error) {
	p, err := a.driver.CreatePost(postModelToDriver(post))
	if err != nil {
		return nil, err
	}
	return postDriverToModel(p), nil
}

func (a *Adapter) RetrievePostByID(i int) (*models.Post, error) {
	p, err := a.driver.RetrievePostByID(i)
	if err != nil {
		return nil, err
	}
	return postDriverToModel(p), nil
}

func (a *Adapter) RetrieveMediaByPostID(id int) ([]*models.Media, error) {
	m, err := a.driver.RetrieveMediaByPostID(id)
	if err != nil {
		return nil, err
	}
	media := make([]*models.Media, len(m))
	for k, v := range m {
		media[k] = mediaDriverToModel(v)
	}
	return media, err
}
