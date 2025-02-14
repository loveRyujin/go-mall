package errcode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
)

type AppError struct {
	code     int    `json:"code"`
	message  string `json:"message"`
	cause    error  `json:"cause"`
	occurred string `json:"occurred"`
}

func (e *AppError) Code() int {
	return e.code
}

func (e *AppError) Message() string {
	return e.message
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	formattedErr := struct {
		Code     int    `json:"code"`
		Message  string `json:"message"`
		Cause    string `json:"cause"`
		Occurred string `json:"occurred"`
	}{
		Code:     e.Code(),
		Message:  e.Message(),
		Occurred: e.occurred,
	}
	if e.cause != nil {
		formattedErr.Cause = e.cause.Error()
	}
	jsonData, _ := json.Marshal(formattedErr)
	return string(jsonData)
}

func (e *AppError) String() string {
	return e.Error()
}

func (e *AppError) HttpStatusCode() int {
	switch e.Code() {
	case Success.Code():
		return http.StatusOK
	case ErrServer.Code(), ErrPanic.Code():
		return http.StatusInternalServerError
	case ErrParams.Code(), ErrUserInvalid.Code():
		return http.StatusBadRequest
	case ErrNotFound.Code():
		return http.StatusNotFound
	case ErrTooManyRequests.Code():
		return http.StatusTooManyRequests
	case ErrToken.Code():
		return http.StatusUnauthorized
	case ErrForbidden.Code():
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

func (e *AppError) WithCause(err error) *AppError {
	newErr := e.Clone()
	newErr.cause = err
	newErr.occurred = getAppErrOccurred()
	return newErr
}

func (e *AppError) UnWrap() error {
	return e.cause
}

func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.code == t.code
}

func (e *AppError) Clone() *AppError {
	return &AppError{
		code:     e.code,
		message:  e.message,
		cause:    e.cause,
		occurred: e.occurred,
	}
}

func newError(code int, message string) *AppError {
	return &AppError{
		code:    code,
		message: message,
	}
}

func Wrap(message string, err error) *AppError {
	if err == nil {
		return nil
	}
	return &AppError{
		code:     -1,
		message:  message,
		cause:    err,
		occurred: getAppErrOccurred(),
	}
}

func getAppErrOccurred() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	file = path.Base(file)
	funcName := runtime.FuncForPC(pc).Name()
	triggerInfo := fmt.Sprintf("func: %s, file: %s, line: %d", funcName, file, line)
	return triggerInfo
}
