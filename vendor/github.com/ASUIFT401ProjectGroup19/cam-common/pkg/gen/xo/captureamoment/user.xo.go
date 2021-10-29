package captureamoment

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// User represents a row from 'captureamoment.user'.
type User struct {
	UserID    int    `json:"UserID"`    // UserID
	FirstName string `json:"FirstName"` // FirstName
	LastName  string `json:"LastName"`  // LastName
	Email     string `json:"Email"`     // Email
	Password  string `json:"Password"`  // Password
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted returns true when the User has been marked for deletion from
// the database.
func (u *User) Deleted() bool {
	return u._deleted
}

// Insert inserts the User to the database.
func (u *User) Insert(ctx context.Context, db DB) error {
	switch {
	case u._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case u._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO captureamoment.user (` +
		`FirstName, LastName, Email, Password` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, u.FirstName, u.LastName, u.Email, u.Password)
	res, err := db.ExecContext(ctx, sqlstr, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	u.UserID = int(id)
	// set exists
	u._exists = true
	return nil
}

// Update updates a User in the database.
func (u *User) Update(ctx context.Context, db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case u._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE captureamoment.user SET ` +
		`FirstName = ?, LastName = ?, Email = ?, Password = ? ` +
		`WHERE UserID = ?`
	// run
	logf(sqlstr, u.FirstName, u.LastName, u.Email, u.Password, u.UserID)
	if _, err := db.ExecContext(ctx, sqlstr, u.FirstName, u.LastName, u.Email, u.Password, u.UserID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the User to the database.
func (u *User) Save(ctx context.Context, db DB) error {
	if u.Exists() {
		return u.Update(ctx, db)
	}
	return u.Insert(ctx, db)
}

// Upsert performs an upsert for User.
func (u *User) Upsert(ctx context.Context, db DB) error {
	switch {
	case u._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO captureamoment.user (` +
		`UserID, FirstName, LastName, Email, Password` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`FirstName = VALUES(FirstName), LastName = VALUES(LastName), Email = VALUES(Email), Password = VALUES(Password)`
	// run
	logf(sqlstr, u.UserID, u.FirstName, u.LastName, u.Email, u.Password)
	if _, err := db.ExecContext(ctx, sqlstr, u.UserID, u.FirstName, u.LastName, u.Email, u.Password); err != nil {
		return logerror(err)
	}
	// set exists
	u._exists = true
	return nil
}

// Delete deletes the User from the database.
func (u *User) Delete(ctx context.Context, db DB) error {
	switch {
	case !u._exists: // doesn't exist
		return nil
	case u._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM captureamoment.user ` +
		`WHERE UserID = ?`
	// run
	logf(sqlstr, u.UserID)
	if _, err := db.ExecContext(ctx, sqlstr, u.UserID); err != nil {
		return logerror(err)
	}
	// set deleted
	u._deleted = true
	return nil
}

// UserByEmail retrieves a row from 'captureamoment.user' as a User.
//
// Generated from index 'user_Email_uindex'.
func UserByEmail(ctx context.Context, db DB, email string) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`UserID, FirstName, LastName, Email, Password ` +
		`FROM captureamoment.user ` +
		`WHERE Email = ?`
	// run
	logf(sqlstr, email)
	u := User{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, email).Scan(&u.UserID, &u.FirstName, &u.LastName, &u.Email, &u.Password); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}

// UserByUserID retrieves a row from 'captureamoment.user' as a User.
//
// Generated from index 'user_UserID_pkey'.
func UserByUserID(ctx context.Context, db DB, userID int) (*User, error) {
	// query
	const sqlstr = `SELECT ` +
		`UserID, FirstName, LastName, Email, Password ` +
		`FROM captureamoment.user ` +
		`WHERE UserID = ?`
	// run
	logf(sqlstr, userID)
	u := User{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, userID).Scan(&u.UserID, &u.FirstName, &u.LastName, &u.Email, &u.Password); err != nil {
		return nil, logerror(err)
	}
	return &u, nil
}