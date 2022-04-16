package cam

import (
	"context"
	"database/sql"
	"encoding/base64"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type Media struct {
	camXO.Media
}

func MediaFromModel(m *types.Media) *Media {
	link := sql.NullString{
		String: base64.StdEncoding.EncodeToString([]byte(m.Link)),
	}
	if m.Link != "" {
		link.Valid = true
	}
	return &Media{
		Media: camXO.Media{
			MediaID:   m.ID,
			MediaLink: link,
			PostID:    m.PostID,
		},
	}
}

func (m *Media) ToModel() *types.Media {
	return &types.Media{
		ID:     m.MediaID,
		Link:   func() string { s, _ := base64.StdEncoding.DecodeString(m.MediaLink.String); return string(s) }(),
		PostID: m.PostID,
	}
}

func (d *Driver) CreateMedia(media *Media) (*Media, error) {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return nil, &BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = media.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return nil, &Unknown{Message: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, &Exists{Message: err.Error()}
		} else {
			return nil, &InsertRecord{Message: err.Error()}
		}
	case nil:
		return media, nil
	}
}

func (d *Driver) RetrieveMediaByPostID(id int) ([]*Media, error) {
	m, err := camXO.MediaByPostID(context.Background(), d.db, id)
	if err != nil {
		return nil, err
	}
	media := make([]*Media, len(m))
	for i, v := range m {
		media[i] = &Media{
			Media: *v,
		}
	}
	return media, nil
}
