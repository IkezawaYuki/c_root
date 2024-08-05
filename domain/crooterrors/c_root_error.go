package crooterrors

import (
	"fmt"
	native "github.com/pkg/errors"
)

type ErrorCode string

const (
	InvalidRequestError ErrorCode = "InvalidRequestError"
	UnauthorizedError   ErrorCode = "UnauthorizedError"
	ForbiddenError      ErrorCode = "ForbiddenError"
	InternalServerError ErrorCode = "InternalServerError"
)

func (code ErrorCode) ToString() string {
	return string(code)
}

type ErrorInfo interface {
	ErrorCode() ErrorCode
	Error() string
}

type errorInfo struct {
	errCode   ErrorCode
	baseError error
}

func New(errorCode ErrorCode, baseError error) ErrorInfo {
	return &errorInfo{
		errCode:   errorCode,
		baseError: baseError,
	}
}

func Cause(err error) error {
	return native.Cause(err)
}

func ExtractError(err error) ErrorInfo {
	var e ErrorInfo
	ok := native.As(Cause(err), &e)
	if !ok {
		return nil
	}
	return e
}

func IsErrorCode(err error, code ErrorCode) bool {
	e := ExtractError(err)
	if e == nil {
		return false
	}
	return e.ErrorCode() == code
}

func Is(err, target error) bool {
	return native.Is(err, target)
}

func (e *errorInfo) Error() string {
	if e.baseError == nil {
		return e.errCode.ToString()
	}
	return fmt.Sprintf("%s %s", e.baseError.Error(), e.baseError.Error())
}

func (e *errorInfo) ErrorCode() ErrorCode {
	return e.errCode
}
