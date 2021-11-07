package cam

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"go.uber.org/zap"
)

func CheckPassword(user *cam.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &ErrorPasswordCheck{msg: err.Error()}
	}
	return nil
}

func (d *Driver) GetUser(username string) (*cam.User, error) {
	u, err := cam.UserByEmail(context.Background(), d.db, username)
	if err != nil {
		return nil, &ErrorUserRetrieval{msg: err.Error()}
	}
	return u, nil
}

func (d *Driver) SetUser(user *cam.User) error {
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

func encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
