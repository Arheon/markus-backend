package auth

import "errors"

const (
	minPassLen = 8
	maxPassLen = 72
)

var (
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must not exided 72 characters")
)

func ValidatePasswordStrangth(passwordString string) error {
	if len(passwordString) < minPassLen {
		return ErrPasswordTooShort
	}

	if len(passwordString) > maxPassLen {
		return ErrPasswordTooLong
	}

	return nil
}
