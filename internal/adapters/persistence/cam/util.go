package cam

import (
	"database/sql"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/models"
)

func userDriverToModel(u *cam.User) *models.User {
	return &models.User{
		ID:        u.UserID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func userModelToDriver(u *models.User) *cam.User {
	return &cam.User{
		UserID:    u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func postDriverToModel(p *cam.Post) *models.Post {
	return &models.Post{
		ID:          p.PostID,
		Description: p.Description.String,
		UserID:      p.UserID,
	}
}

func postModelToDriver(p *models.Post) *cam.Post {
	description := sql.NullString{
		String: p.Description,
	}
	if p.Description != "" {
		description.Valid = true
	}
	return &cam.Post{
		PostID:      p.ID,
		Description: description,
		UserID:      p.UserID,
	}
}

func mediaDriverToModel(m *cam.Media) *models.Media {
	return &models.Media{
		ID:     m.MediaID,
		Link:   m.MediaLink.String,
		PostID: m.Postid,
	}
}

func mediaModelToDriver(m *models.Media) *cam.Media {
	link := sql.NullString{
		String: m.Link,
	}
	if m.Link != "" {
		link.Valid = true
	}
	return &cam.Media{
		MediaID:   m.ID,
		MediaLink: link,
		Postid:    m.PostID,
	}
}
