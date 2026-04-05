package auth

import (
	"errors"
	"strings"
	"unicode"
)

const (
	minUsernameLength = 2
	maxUsernameLength = 32
)

var (
	ErrUsernameTooShort               = errors.New("username must be at least 2 characters")
	ErrUsernameTooLong                = errors.New("username must be at most 32 characters")
	ErrUsernameHasControlCharacters   = errors.New("username must not contain control characters")
	ErrUsernameHasInvisibleCharacters = errors.New("username has invisible characters")
)

func ValidateUsername(username string) error {
	username = strings.TrimSpace(username)
	n := len([]rune(username))

	if n < minUsernameLength {
		return ErrUsernameTooShort
	}

	if n > maxUsernameLength {
		return ErrUsernameTooLong
	}

	for _, r := range username {
		if unicode.IsControl(r) {
			return ErrUsernameHasControlCharacters
		}

		if unicode.In(r, unicode.Cf) {
			return ErrUsernameHasInvisibleCharacters
		}
	}

	return nil
}
