package cam

import (
	"database/sql"
	"encoding/base64"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
)

func userDriverToModel(u *camXO.User) *types.User {
	return &types.User{
		ID:        u.UserID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func userModelToDriver(u *types.User) *camXO.User {
	return &camXO.User{
		UserID:    u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func postDriverToModel(p *camXO.Post) *types.Post {
	return &types.Post{
		ID:          p.PostID,
		Description: p.Description.String,
		UserID:      p.UserID,
		Date:        p.Date,
	}
}

func postModelToDriver(p *types.Post) *camXO.Post {
	description := sql.NullString{
		String: p.Description,
	}
	if p.Description != "" {
		description.Valid = true
	}
	return &camXO.Post{
		PostID:      p.ID,
		Description: description,
		UserID:      p.UserID,
	}
}

func mediaDriverToModel(m *camXO.Media) *types.Media {
	return &types.Media{
		ID:     m.MediaID,
		Link:   func() string { s, _ := base64.StdEncoding.DecodeString(m.MediaLink.String); return string(s) }(),
		PostID: m.PostID,
	}
}

func mediaModelToDriver(m *types.Media) *camXO.Media {
	link := sql.NullString{
		String: base64.StdEncoding.EncodeToString([]byte(m.Link)),
	}
	if m.Link != "" {
		link.Valid = true
	}
	return &camXO.Media{
		MediaID:   m.ID,
		MediaLink: link,
		PostID:    m.PostID,
	}
}
