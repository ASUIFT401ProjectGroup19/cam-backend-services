package cam

import (
	camDriver "github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/adapters/persistence/cam/database/cam"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
)

// Adapter implements feed.Storage.
// Adapter implements identity.Storage.
// Adapter implements post.Storage.
// Adapter implements subscription.Storage.
// Adapter serves to convert between cam.driver types and application domain types in core.types.
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
	return user.ToModel(), nil
}

func (a *Adapter) CreateUser(user *types.User) (*types.User, error) {
	u, err := a.driver.CreateUser(camDriver.UserFromModel(user))
	if err != nil {
		return nil, err
	}
	return u.ToModel(), nil
}

func (a *Adapter) RetrieveUserByID(id int) (*types.User, error) {
	user, err := a.driver.RetrieveUserByID(id)
	if err != nil {
		return nil, err
	}
	return user.ToModel(), nil
}

func (a *Adapter) RetrieveUserByUserName(username string) (*types.User, error) {
	user, err := a.driver.RetrieveUserByUserName(username)
	if err != nil {
		return nil, err
	}
	return user.ToModel(), nil
}

func (a *Adapter) CreateMedia(media *types.Media) (*types.Media, error) {
	m, err := a.driver.CreateMedia(camDriver.MediaFromModel(media))
	if err != nil {
		return nil, err
	}
	return m.ToModel(), nil
}

func (a *Adapter) CreatePost(post *types.Post) (*types.Post, error) {
	p, err := a.driver.CreatePost(camDriver.PostFromModel(post))
	if err != nil {
		return nil, err
	}
	return p.ToModel(), nil
}

func (a *Adapter) RetrievePostByID(i int) (*types.Post, error) {
	p, err := a.driver.RetrievePostByID(i)
	if err != nil {
		return nil, err
	}
	return p.ToModel(), nil
}

func (a *Adapter) RetrieveSubscribedPostsPaginated(userID, pageNumber, batchSize int) ([]*types.Post, error) {
	p, err := a.driver.RetrieveSubscribedPostsPaginated(userID, pageNumber, batchSize)
	if err != nil {
		return nil, err
	}
	posts := make([]*types.Post, len(p))
	for i, v := range p {
		posts[i] = v.ToModel()
	}
	return posts, nil
}

func (a *Adapter) RetrieveUserPostsPaginated(userID, pageNumber, batchSize int) ([]*types.Post, error) {
	p, err := a.driver.RetrieveUserPostsPaginated(userID, pageNumber, batchSize)
	if err != nil {
		return nil, err
	}
	posts := make([]*types.Post, len(p))
	for i, v := range p {
		posts[i] = v.ToModel()
	}
	return posts, nil
}

func (a *Adapter) RetrieveMediaByPostID(id int) ([]*types.Media, error) {
	m, err := a.driver.RetrieveMediaByPostID(id)
	if err != nil {
		return nil, err
	}
	media := make([]*types.Media, len(m))
	for k, v := range m {
		media[k] = v.ToModel()
	}
	return media, err
}

func (a *Adapter) CreateSubscription(user, other *types.User) error {
	err := a.driver.CreateSubscription(camDriver.SubFromModel(&types.Sub{UserID: user.ID, OtherID: other.ID}))
	if err != nil {
		return err
	}
	return nil
}

func (a *Adapter) DeleteSubscription(user, other *types.User) error {
	err := a.driver.DeleteSubscription(camDriver.SubFromModel(&types.Sub{UserID: user.ID, OtherID: other.ID}))
	if err != nil {
		return err
	}
	return nil
}
