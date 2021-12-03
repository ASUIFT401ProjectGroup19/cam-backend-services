package cam

import (
	"context"
	"database/sql"

	cam "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/errors"
)

func (d Driver) CheckPassword(user *cam.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &errors.PasswordCheck{Message: err.Error()}
	}
	return nil
}

func (d *Driver) CreateUser(user *cam.User) (*cam.User, error) {
	hash, err := encrypt(user.Password)
	if err != nil {
		return nil, &errors.EncryptPassword{Message: err.Error()}
	}
	user.Password = hash
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
	err = user.Insert(context.Background(), d.db)
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
		return user, nil
	}
}

func (d *Driver) RetrieveUserByUserName(username string) (*cam.User, error) {
	u, err := cam.UserByEmail(context.Background(), d.db, username)
	if err != nil {
		return nil, &errors.UserRetrieval{Message: err.Error()}
	}
	return u, nil
}

// func UpdateUser
// func DeleteUser

func encrypt(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
