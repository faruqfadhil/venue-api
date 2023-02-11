package error

import (
	"errors"
	"fmt"
)

var (
	ErrGeneralBadRequest = errors.New("bad request")
	ErrGeneralNotFound   = errors.New("not found")
	ErrGeneralDB         = errors.New("DB error")
	ErrInternal          = errors.New("internal server error")
)

type InternalError struct {
	UserErrMsg  string
	OriginalErr error
	TypeErr     error
}

func New(typeErr error, originalErr error, userErrMsg ...string) error {
	var msg string
	if len(userErrMsg) > 0 {
		msg = userErrMsg[0]
	}
	return &InternalError{
		UserErrMsg:  msg,
		OriginalErr: originalErr,
		TypeErr:     typeErr,
	}
}

func (d *InternalError) Error() string {
	if d.UserErrMsg == "" {
		return "unexpected error"
	}
	return d.UserErrMsg
}

func toInternalErr(err error) *InternalError {
	var intErr *InternalError
	e, ok := err.(*InternalError)
	if ok {
		intErr = e
	} else {
		intErr = &InternalError{
			UserErrMsg:  "internal service error",
			OriginalErr: fmt.Errorf("unexpected err: %v", err),
			TypeErr:     ErrInternal,
		}
	}
	return intErr
}

func GetTypeErr(err error) error {
	return toInternalErr(err).TypeErr
}
