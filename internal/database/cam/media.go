package cam

import (
	"context"
	"database/sql"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errs"
)

func (d *Driver) CreateMedia(media *cam.Media) (*cam.Media, error) {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return nil, &errs.BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = media.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return nil, &errs.Unknown{Message: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, &errs.Exists{Message: err.Error()}
		} else {
			return nil, &errs.InsertRecord{Message: err.Error()}
		}
	case nil:
		return media, nil
	}
}

func (d *Driver) RetrieveMediaByPostID(id int) ([]*cam.Media, error) {
	m, err := cam.MediaByPostid(context.Background(), d.db, id)
	if err != nil {
		return nil, err
	}
	return m, nil
}
