package cam

import (
	"context"
	"database/sql"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errors"
)

func (d *Driver) CreatePost(post *cam.Post) (*cam.Post, error) {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return nil, &errors.BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = post.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return nil, &errors.Unknown{Message: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return nil, &errors.Exists{Message: err.Error()}
		} else {
			return nil, &errors.InsertRecord{Message: err.Error()}
		}
	case nil:
		return post, nil
	}
}
