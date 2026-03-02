package apperror

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

type ValidationError struct {
	Param   string `json:"param"`
	Message string `json:"message"`
}

func NewValidationError(param string, message string) ValidationError {
	return ValidationError{
		Param:   param,
		Message: message,
	}
}

func MsgForTag(fe validator.FieldError) string {
	field := strings.ToLower(fe.Field())

	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("Поле '%s' обязательно", field)
	case "email":
		return fmt.Sprintf("Поле '%s' не соответствует email формату", field)
	case "uuid":
		return fmt.Sprintf("Поле '%s' не соответствует uuid формату", field)
	case "date_by_format":
		return fmt.Sprintf("Поле '%s' не соответствует формату %s", field, fe.Param())
	case "username":
		return "Не верное имя"
	case "phone_format", "numeric_phone_format":
		return fmt.Sprintf("Поле '%s' не соответствует формату", field)
	case "user_status":
		return fmt.Sprintf("Поле '%s' имеет не верное значение", field)
	case "date_range":
		return fmt.Sprintf("Поле '%s' должно принимать значение created_at или updated_at", field)
	case "sort_type":
		return fmt.Sprintf("Поле '%s' должно принимать значение asc или desc", field)
	}
	return fe.Error()
}

func BeautifyErrors(ve validator.ValidationErrors) []ValidationError {
	out := make([]ValidationError, len(ve))
	for i, fe := range ve {
		out[i] = ValidationError{strings.ToLower(fe.Field()), MsgForTag(fe)}
	}

	return out
}
