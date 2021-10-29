package camdb

import "fmt"

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

type ErrorUnknown struct {
	msg string
}

func (e *ErrorUnknown) Error() string {
	return fmt.Sprintf("unknown: %s", e.msg)
}
