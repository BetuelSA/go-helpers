package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

// ErrorType is the type of an error
type ErrorType uint

const (
	// NoType -
	NoType ErrorType = iota
	// BadRequest (400)
	BadRequest
	// Unauthorized (401)
	Unauthorized
	// Forbidden (403)
	Forbidden
	// NotFound (404)
	NotFound
	// MethodNotAllowed (405)
	MethodNotAllowed
	// PreconditionFailed (412)
	PreconditionFailed
	// UnsupportedMediaType (415)
	UnsupportedMediaType
	// InternalServerError (500)
	InternalServerError
	// NotImplemented (501)
	NotImplemented
	// ServiceUnavailable (503)
	ServiceUnavailable
)

type customError struct {
	errorType ErrorType
	message   error
	detail    error
	context   errorContext
}

type errorContext struct {
	Field   string
	Message string
}

// New creates a new customError
func (errorType ErrorType) New(msg string) error {
	return errorType.Newf(msg)
}

// Newf creates a new customError with formatted message
func (errorType ErrorType) Newf(msg string, args ...interface{}) error {
	errMsg := fmt.Errorf(msg, args...)
	errDetail := fmt.Errorf("")
	return customError{errorType: errorType, message: errMsg, detail: errDetail}
}

// Wrap creates a new wrapped error
func (errorType ErrorType) Wrap(err error, msg string) error {
	return errorType.Wrapf(err, msg)
}

// Wrapf creates a new wrapped error with formatted message
func (errorType ErrorType) Wrapf(err error, msg string, args ...interface{}) error {
	errMsg := fmt.Errorf(msg, args...)
	errDetail := errors.Wrap(err, errMsg.Error())
	return customError{errorType: errorType, message: errMsg, detail: errDetail}
}

// Error returns the message of a customError
func (error customError) Error() string {
	return error.detail.Error()
}

// New creates a no type error
func New(msg string) error {
	return Newf(msg)
}

// Newf creates a no type error with formatted message
func Newf(msg string, args ...interface{}) error {
	errMsg := fmt.Errorf(msg, args...)
	errDetail := fmt.Errorf("")
	return customError{errorType: NoType, message: errMsg, detail: errDetail}
}

// Wrap an error with a string
func Wrap(err error, msg string) error {
	return Wrapf(err, msg)
}

// Wrapf an error with format string
func Wrapf(err error, msg string, args ...interface{}) error {
	errMsg := fmt.Errorf(msg, args...)
	errDetail := errors.Wrap(err, errMsg.Error())

	if customErr, ok := err.(customError); ok {
		return customError{
			errorType: customErr.errorType,
			message:   errMsg,
			detail:    errDetail,
			context:   customErr.context,
		}
	}

	return customError{errorType: NoType, message: errMsg, detail: errDetail}
}

// Cause gives the original error
func Cause(err error) error {
	return errors.Cause(err)
}

// GetErrorMessage return the detail of a customError
func GetErrorMessage(err error) string {
	if customErr, ok := err.(customError); ok {
		return customErr.message.Error()
	}

	return err.Error()
}

// GetErrorDetail return the detail of a customError
func GetErrorDetail(err error) string {
	if customErr, ok := err.(customError); ok {
		return customErr.detail.Error()
	}

	return ""
}

// AddErrorContext adds a context to an error
func AddErrorContext(err error, field, message string) error {
	context := errorContext{Field: field, Message: message}
	if customErr, ok := err.(customError); ok {
		return customError{errorType: customErr.errorType, message: customErr.message, detail: customErr.detail, context: context}
	}

	return customError{errorType: NoType, message: err, context: context}
}

// GetErrorContext returns the error context
func GetErrorContext(err error) map[string]string {
	emptyContext := errorContext{}
	if customErr, ok := err.(customError); ok && customErr.context != emptyContext {

		return map[string]string{"field": customErr.context.Field, "message": customErr.context.Message}
	}

	return nil
}

// GetType returns the error type
func GetType(err error) ErrorType {
	if customErr, ok := err.(customError); ok {
		return customErr.errorType
	}

	return NoType
}
