package camdb

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"go.uber.org/zap"
)

func (d *DB) SetUser(user *cam.User) error {
	hash, err := encrypt(user.Password)
	if err != nil {
		return &ErrorEncryptPassword{msg: err.Error()}
	}
	user.Password = hash
	transaction, err := d.db.Beginx()
	if err != nil {
		d.log.Error("database begin transaction", zap.Error(err))
		return &ErrorBeginTransaction{msg: err.Error()}
	}
	defer func() {
		if err := transaction.Rollback(); err != nil && err != sql.ErrTxDone {
			d.log.Error("database error; rollback transaction", zap.Error(err))
		}
	}()
	err = user.Insert(context.Background(), d.db)
	switch err.(type) {
	default:
		return &ErrorUnknown{msg: err.Error()}
	case *mysql.MySQLError:
		if err.(*mysql.MySQLError).Number == 1062 {
			return &ErrorExists{msg: err.Error()}
		} else {
			return &ErrorInsertRecord{msg: err.Error()}
		}
	case nil:
		return nil
	}
}

func encrypt(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
