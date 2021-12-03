package errors

import (
	"fmt"
)

type BeginTransaction struct {
	Message string
}

func (e *BeginTransaction) Error() string {
	return fmt.Sprintf("begin transaction: %s", e.Message)
}

type EncryptPassword struct {
	Message string
}

func (e *EncryptPassword) Error() string {
	return fmt.Sprintf("encrypt password: %s", e.Message)
}

type Exists struct {
	Message string
}

func (e *Exists) Error() string {
	return fmt.Sprintf("record exists: %s", e.Message)
}

type InsertRecord struct {
	Message string
}

func (e *InsertRecord) Error() string {
	return fmt.Sprintf("insert record: %s", e.Message)
}

type PasswordCheck struct {
	Message string
}

func (e *PasswordCheck) Error() string {
	return fmt.Sprintf("password check: %s", e.Message)
}

type Unknown struct {
	Message string
}

func (e *Unknown) Error() string {
	return fmt.Sprintf("unknown: %s", e.Message)
}

type UnsupportedDriver struct {
	Message string
}

func (e *UnsupportedDriver) Error() string {
	return fmt.Sprintf("unsupported driver: %s", e.Message)
}

type UserRetrieval struct {
	Message string
}

func (e *UserRetrieval) Error() string {
	return fmt.Sprintf("user not found: %s", e.Message)
}
