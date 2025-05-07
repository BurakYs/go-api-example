package middleware

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	"github.com/BurakYs/GoAPIExample/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateBody[T any]() gin.HandlerFunc {
	return validate[T]("body")
}

func ValidateQuery[T any]() gin.HandlerFunc {
	return validate[T]("query")
}

func ValidateParams[T any]() gin.HandlerFunc {
	return validate[T]("uri")
}

type normalizable interface {
	Normalize()
}

func validate[T any](kind string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var data T
		var err error

		switch kind {
		case "body":
			err = ctx.ShouldBindJSON(&data)
		case "query":
			err = ctx.ShouldBindQuery(&data)
		case "uri":
			err = ctx.ShouldBindUri(&data)
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, formatValidationError(err, data))
			return
		}

		if normalizable, ok := any(&data).(normalizable); ok {
			normalizable.Normalize()
		}

		ctx.Set(kind, data)
		ctx.Next()
	}
}

func formatValidationError(err error, data any) models.ValidationError {
	var ve validator.ValidationErrors

	switch {
	case errors.As(err, &ve):
		return formatValidatorErrors(&ve, data)
	case err.Error() == "EOF":
		return models.ValidationError{
			Message:            "Invalid JSON body",
			ValidationFailures: []models.ValidationFailure{},
		}
	default:
		return models.ValidationError{
			Message:            "Invalid parameters provided",
			ValidationFailures: []models.ValidationFailure{},
		}
	}
}

func formatValidatorErrors(ve *validator.ValidationErrors, data any) models.ValidationError {
	failures := make([]models.ValidationFailure, 0, len(*ve))

	for _, fe := range *ve {
		field := getJSONFieldName(fe.StructField(), data)
		failures = append(failures, formatFieldError(fe, field))
	}

	return models.ValidationError{
		Message:            "Invalid parameters provided",
		ValidationFailures: failures,
	}
}

func formatFieldError(fe validator.FieldError, field string) models.ValidationFailure {
	switch fe.Tag() {
	case "required":
		return models.ValidationFailure{
			Type:    "required",
			Field:   field,
			Message: "This field is required",
		}
	case "min":
		return models.ValidationFailure{
			Type:    "min",
			Field:   field,
			Message: "This field must be at least " + fe.Param() + " characters",
		}
	case "max":
		return models.ValidationFailure{
			Type:    "max",
			Field:   field,
			Message: "This field must be at most " + fe.Param() + " characters",
		}
	default:
		return models.ValidationFailure{
			Type:    "invalid",
			Field:   field,
			Message: "This field is invalid",
		}
	}
}

func getJSONFieldName(structField string, obj any) string {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(structField); ok {
		tag := f.Tag.Get("json")

		if tag == "-" {
			return ""
		}

		if tag != "" {
			return strings.Split(tag, ",")[0]
		}
	}

	return ""
}
