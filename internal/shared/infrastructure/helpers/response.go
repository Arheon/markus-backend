package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	domainerrors "github.com/Arheon/markus-backend/internal/shared/domain/errors"
	"github.com/Arheon/markus-backend/internal/shared/infrastructure/apperror"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const defValidationSubtitle = "Проверь наличие опечаток, пробелов и спецсимволов"

func AbortWithDomainErrorJSON(c *gin.Context, e domainerrors.Error) {
	AbortWithError400JSON(
		c,
		e.Error(),
		int(e.Code()),
		e.Title(),
		e.Subtitle(),
	)
}

func AbortWithError400JSON(c *gin.Context, message string, code int, title, subtitle string) {
	c.AbortWithStatusJSON(
		http.StatusBadRequest,
		gin.H{
			"status":   "error",
			"techwork": nil,
			"error": map[string]any{
				"code":    code,
				"message": message,
				"message_ext": map[string]any{
					"title":    title,
					"subtitle": subtitle,
				},
			},
		},
	)
}

func AbortWithError400ECJSON(c *gin.Context, code domainerrors.ErrorCode, title, subtitle string) {
	AbortWithError400JSON(
		c,
		"",
		int(code),
		title,
		subtitle,
	)
}

// AbortWithError400ValidationJSON для обработки ошибок валидации в случаях, когда
// ошибки обрабатываются через internal/servers/http/middlewares/error_handler.go
func AbortWithError400ValidationJSON(c *gin.Context, err error) {
	var (
		ve validator.ValidationErrors
		ue *json.UnmarshalTypeError
		ne *strconv.NumError
	)

	switch {
	case errors.As(err, &ve):
		isSetPhone := slices.ContainsFunc(ve, func(fieldError validator.FieldError) bool {
			return strings.ToLower(fieldError.Field()) == "phone"
		})

		if isSetPhone {
			AbortWithError400JSON(
				c,
				"",
				int(domainerrors.EC15),
				"Неверный формат телефона",
				"Попробуй ввести номер телефона повторно",
			)
			return
		}

		firstVe := apperror.BeautifyErrors(ve)[0]
		AbortWithError400JSON(c, "", int(domainerrors.EC31), firstVe.Message, defValidationSubtitle)
		return
	case errors.As(err, &ue):
		errorMessage := fmt.Sprintf("Значение '%s' не соответствует типу данных", ue.Field)
		AbortWithError400JSON(c, "", int(domainerrors.EC31), errorMessage, defValidationSubtitle)
		return
	case errors.As(err, &ne):
		errorMessage := fmt.Sprintf("Значение '%s' не соответствует формату", ne.Num)
		AbortWithError400JSON(c, "", int(domainerrors.EC31), errorMessage, defValidationSubtitle)
		return
	}

	AbortWithInternalErrorJSON(c)
}

func AbortWithStatusJSON(c *gin.Context, statusCode int, message string, errorCode int) {
	c.AbortWithStatusJSON(
		statusCode,
		gin.H{
			"status": "error",
			"error": map[string]any{
				"code":    errorCode,
				"message": message,
			},
		},
	)
}

func AbortWithStatusValidatesJSON(c *gin.Context, statusCode int, err any, errorCode int) {
	c.AbortWithStatusJSON(
		statusCode,
		gin.H{
			"status": "error",
			"error": map[string]any{
				"code":            errorCode,
				"validate_errors": err,
			},
		},
	)
}

func AbortWithError(c *gin.Context, err error) {
	c.Header("Content-Type", "application/json")
	c.AbortWithError(http.StatusBadRequest, err) //nolint:all
}

func SuccessJSONV2(c *gin.Context, response any) {
	c.AbortWithStatusJSON(http.StatusOK, gin.H{
		"status":   "ok",
		"response": response,
	})
}

func AbortWithRequiredParameterJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusBadRequest, "Required parameter", 460)
}

func AbortWithNotFoundJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusBadRequest, "Not Found", 404)
}

func AbortWithInternalErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusInternalServerError, "Server internal error", 500)
}

func AbortWithBadRequestErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusBadRequest, "Bad Request", 400)
}

func AbortWithEntityExistErrorMessageJSON(c *gin.Context, message string) {
	AbortWithStatusJSON(c, 400, message, 463)
}

// Deprecated: AbortWithValidationErrorJSON is deprecated. Use AbortWithError400ValidationJSON for error handler
func AbortWithValidationErrorJSON(c *gin.Context, err error) {
	var (
		ve validator.ValidationErrors
		ue *json.UnmarshalTypeError
		ne *strconv.NumError
	)

	switch {
	case errors.As(err, &ve):
		AbortWithStatusValidatesJSON(c, http.StatusBadRequest, apperror.BeautifyErrors(ve), int(domainerrors.EC31))
		return
	case errors.As(err, &ue):
		errorMessage := fmt.Sprintf("Значение '%s' не соответствует типу данных", ue.Field)
		formatErrors := []apperror.ValidationError{apperror.NewValidationError(ue.Field, errorMessage)}

		AbortWithStatusValidatesJSON(c, http.StatusBadRequest, formatErrors, int(domainerrors.EC31))
		return
	case errors.As(err, &ne):
		errorMessage := fmt.Sprintf("Значение '%s' не соответствует формату", ne.Num)
		formatErrors := []apperror.ValidationError{apperror.NewValidationError("", errorMessage)}

		AbortWithStatusValidatesJSON(c, http.StatusBadRequest, formatErrors, int(domainerrors.EC31))
		return
	}

	AbortWithInternalErrorJSON(c)
}

func AbortWithAccessDeniedErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusBadRequest, "Bad Request", 470)
}

func AbortWithUnauthorizedErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusBadRequest, "Unauthorized", 462) //nolint:gomnd
}

func AbortWithFailedDependencyErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusFailedDependency, "Failed Dependency", 424) //nolint:gomnd
}
func AbortWithUnprocessableEntityErrorJSON(c *gin.Context) {
	AbortWithStatusJSON(c, http.StatusUnprocessableEntity, "Unprocessable Entity", 422) //nolint:gomnd
}
