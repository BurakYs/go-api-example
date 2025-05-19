package middleware

import (
	"cmp"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/BurakYs/GoAPIExample/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type BindingLocation string

const (
	BindingLocationBody          BindingLocation = "body"
	BindingLocationQuery         BindingLocation = "query"
	BindingLocationParams        BindingLocation = "params"
	BindingLocationForm          BindingLocation = "form"
	BindingLocationMultipartForm BindingLocation = "multipartForm"
)

func ValidateBody[T any]() gin.HandlerFunc {
	return validate[T](BindingLocationBody)
}

func ValidateQuery[T any]() gin.HandlerFunc {
	return validate[T](BindingLocationQuery)
}

func ValidateParams[T any]() gin.HandlerFunc {
	return validate[T](BindingLocationParams)
}

func ValidateForm[T any]() gin.HandlerFunc {
	return validate[T](BindingLocationForm)
}

func ValidateMultipartForm[T any]() gin.HandlerFunc {
	return validate[T](BindingLocationMultipartForm)
}

type transformable interface {
	Transform()
}

var vld = validator.New()

func validate[T any](location BindingLocation) gin.HandlerFunc {
	return func(c *gin.Context) {
		var data T
		var err error

		switch location {
		case BindingLocationBody:
			err = c.ShouldBindJSON(&data)
		case BindingLocationQuery:
			err = c.ShouldBindQuery(&data)
		case BindingLocationParams:
			err = c.ShouldBindUri(&data)
		case BindingLocationForm:
			err = c.ShouldBindWith(&data, binding.Form)
		case BindingLocationMultipartForm:
			err = c.ShouldBindWith(&data, binding.FormMultipart)
		}

		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, formatValidationError(err, string(location), data))
			return
		}

		if err := vld.Struct(data); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, formatValidationError(err, string(location), data))
			return
		}

		if t, ok := any(&data).(transformable); ok {
			t.Transform()
		}

		c.Set(string(location), data)
		c.Next()
	}
}

func formatValidationError(err error, location string, data any) any {
	var ve validator.ValidationErrors

	switch {
	case errors.As(err, &ve):
		return formatValidatorErrors(&ve, location, data)
	case err.Error() == "EOF":
		return models.APIError{
			Message: "Empty request body",
		}
	default:
		return models.APIError{
			Message: "Invalid parameters provided",
		}
	}
}

func formatValidatorErrors(ve *validator.ValidationErrors, location string, data any) models.ValidationError {
	failures := make([]models.ValidationFailure, 0, len(*ve))

	for _, fe := range *ve {
		field := getFieldName(fe.StructField(), data)
		failures = append(failures, formatFieldError(fe, location, field))
	}

	return models.ValidationError{
		Message:            "Invalid parameters provided",
		ValidationFailures: failures,
	}
}

func formatFieldError(fe validator.FieldError, location string, field string) models.ValidationFailure {
	msg := "This field is invalid"

	switch fe.Tag() {
	case "required":
		msg = "This field is required"
	case "email":
		msg = "This field must be a valid email address"
	case "min":
		msg = fmt.Sprintf("This field must be at least %s characters", fe.Param())
	case "max":
		msg = fmt.Sprintf("This field must be at most %s characters", fe.Param())
	default:
		msg = fmt.Sprintf("This field is invalid: %s", fe.Tag())
	}

	return models.ValidationFailure{
		Location: location,
		Field:    field,
		Message:  msg,
	}
}

func getFieldName(structField string, obj any) string {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if f, ok := t.FieldByName(structField); ok {
		tag := cmp.Or(f.Tag.Get("json"), f.Tag.Get("uri"), f.Tag.Get("form"))

		if tag == "-" {
			return ""
		}

		if tag != "" {
			return strings.Split(tag, ",")[0]
		}
	}

	return ""
}
