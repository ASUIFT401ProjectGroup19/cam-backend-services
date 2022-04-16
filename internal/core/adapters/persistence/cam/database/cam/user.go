package cam

import (
	"context"
	"database/sql"
	"github.com/ASUIFT401ProjectGroup19/cam-backend-services/internal/core/types"

	camXO "github.com/ASUIFT401ProjectGroup19/cam-common/pkg/gen/xo/captureamoment"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	camXO.User
}

func UserFromModel(u *types.User) *User {
	return &User{
		User: camXO.User{
			UserID:    u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Email:     u.Email,
			Password:  u.Password,
		},
	}
}

func (u *User) ToModel() *types.User {
	return &types.User{
		ID:        u.UserID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
		Password:  u.Password,
	}
}

func (d Driver) CheckPassword(user *User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return &PasswordCheck{Message: err.Error()}
	}
	return nil
}

func (d *Driver) CreateUser(user *User) (*User, error) {
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

func (d *Driver) RetrieveUserByID(id int) (*User, error) {
	u, err := camXO.UserByUserID(context.Background(), d.db, id)
	if err != nil {
		return nil, &UserRetrieval{Message: err.Error()}
	}
	return &User{
		User: *u,
	}, nil
}

func (d *Driver) RetrieveUserByUserName(username string) (*User, error) {
	u, err := camXO.UserByEmail(context.Background(), d.db, username)
	if err != nil {
		return nil, &UserRetrieval{Message: err.Error()}
	}
	return &User{
		User: *u,
	}, nil
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
