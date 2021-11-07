package cam

import (
	"fmt"
)

type ErrorBeginTransaction struct {
	msg string
}

func (e *ErrorBeginTransaction) Error() string {
	return fmt.Sprintf("begin transaction: %s", e.msg)
}

type ErrorEncryptPassword struct {
	msg string
}

func (e *ErrorEncryptPassword) Error() string {
	return fmt.Sprintf("encrypt password: %s", e.msg)
}

type ErrorExists struct {
	msg string
}

func (e *ErrorExists) Error() string {
	return fmt.Sprintf("record exists: %s", e.msg)
}

type ErrorInsertRecord struct {
	msg string
}

func (e *ErrorInsertRecord) Error() string {
	return fmt.Sprintf("insert record: %s", e.msg)
}

type ErrorPasswordCheck struct {
	msg string
}

func (e *ErrorPasswordCheck) Error() string {
	return fmt.Sprintf("password check: %s", e.msg)
}

type ErrorUnknown struct {
	msg string
}

func (e *ErrorUnknown) Error() string {
	return fmt.Sprintf("unknown: %s", e.msg)
}

type ErrorUnsupportedDriver struct {
	msg string
}

func (e *ErrorUnsupportedDriver) Error() string {
	return fmt.Sprintf("unsupported driver: %s", e.msg)
}

type ErrorUserRetrieval struct {
	msg string
}

func (e *ErrorUserRetrieval) Error() string {
	return fmt.Sprintf("user not found: %s", e.msg)
}
