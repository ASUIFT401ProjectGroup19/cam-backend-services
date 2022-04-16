package cam

import (
	"context"
	"database/sql"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"
	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
)

type Sub struct {
	camXO.Subscription
}

func SubFromModel(sub *types.Sub) *Sub {
	return &Sub{
		Subscription: camXO.Subscription{
			UserID:         sub.UserID,
			FollowedUserID: sub.OtherID,
		},
	}
}

func (d *Driver) CreateSubscription(sub *Sub) error {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return &BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = sub.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return &Unknown{Message: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return &Exists{Message: err.Error()}
		} else {
			return &InsertRecord{Message: err.Error()}
		}
	case nil:
		return nil
	}
}

func (d *Driver) DeleteSubscription(sub *Sub) error {
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return &BeginTransaction{Message: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = sub.Delete(context.Background(), d.db)
	switch err {
	default:
		return &Unknown{Message: err.Error()}
	case nil:
		return nil
	}
}
