package cam

import (
	"context"
	"database/sql"
	"time"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

func (d *Driver) CreatePost(post *camXO.Post) (*camXO.Post, error) {
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
	post.Date = time.Now()
	err = post.Insert(context.Background(), d.db)
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
		return post, nil
	}
}

func (d Driver) RetrievePostByID(id int) (*camXO.Post, error) {
	p, err := camXO.PostByPostID(context.Background(), d.db, id)
	if err != nil {
		return nil, err
	}
	return p, nil
}
