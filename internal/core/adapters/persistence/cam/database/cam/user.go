package cam

import (
	"context"
	"database/sql"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (d Driver) CheckPassword(user *camXO.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &PasswordCheck{Message: err.Error()}
	}
	return nil
}

func (d *Driver) CreateUser(user *camXO.User) (*camXO.User, error) {
	hash, err := encrypt(user.Password)
	if err != nil {
		return nil, &EncryptPassword{Message: err.Error()}
	}
	user.Password = hash
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
	err = user.Insert(context.Background(), d.db)
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
		return user, nil
	}
}

func (d *Driver) RetrieveUserByID(id int) (*camXO.User, error) {
	u, err := camXO.UserByUserID(context.Background(), d.db, id)
	if err != nil {
		return nil, &UserRetrieval{Message: err.Error()}
	}
	return u, nil
}

func (d *Driver) RetrieveUserByUserName(username string) (*camXO.User, error) {
	u, err := camXO.UserByEmail(context.Background(), d.db, username)
	if err != nil {
		return nil, &UserRetrieval{Message: err.Error()}
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
