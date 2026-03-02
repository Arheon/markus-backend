package errors

import "errors"

type ErrorCode int

const (
	EC0   ErrorCode = 0
	EC15  ErrorCode = 15
	EC20  ErrorCode = 20
	EC21  ErrorCode = 21
	EC22  ErrorCode = 22
	EC23  ErrorCode = 23
	EC24  ErrorCode = 24
	EC25  ErrorCode = 25
	EC28  ErrorCode = 28 // Ошибка во внешних зависимостях (сервисах)
	EC29  ErrorCode = 29
	EC30  ErrorCode = 30
	EC31  ErrorCode = 31 // Ошибка валидации входящих данных
	EC404 ErrorCode = 404
	EC400 ErrorCode = 400
)

type Error struct {
	message  string
	code     ErrorCode
	title    string
	subtitle string
}

func (e Error) Error() string {
	return e.message
}

func (e Error) SetError(err error) Error {
	e.message = err.Error()

	return e
}

func (e Error) Code() ErrorCode {
	return e.code
}

func (e Error) Title() string {
	return e.title
}

func (e Error) Subtitle() string {
	return e.subtitle
}

func NewError(code ErrorCode, title, subtitle string) Error {
	return Error{
		code:     code,
		title:    title,
		subtitle: subtitle,
	}
}

var ErrDomain = errors.New("domain error")
var ErrInfrastructure = errors.New("infrastructure error")

// Cause выдает оригинальную ошибку, с которой все началось
func Cause(err error) error {
	switch x := err.(type) {
	case interface{ Unwrap() error }:
		e := x.Unwrap()
		if e == nil {
			return err
		}

		return e
	case interface{ Unwrap() []error }:
		s := x.Unwrap()
		if s == nil {
			return nil
		}

		e := s[len(s)-1]
		if e == nil {
			return err
		}
		return e
	default:
		return err
	}
}
