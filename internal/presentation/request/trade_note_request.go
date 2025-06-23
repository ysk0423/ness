package request

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrorResponse struct {
	Errors []ValidationError `json:"errors"`
}

var validate = validator.New()

func ValidateRequest(c echo.Context, req interface{}) error {
	if err := validate.Struct(req); err != nil {
		var errors []ValidationError
		
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Message: getValidationMessage(err),
			})
		}

		return c.JSON(http.StatusBadRequest, ValidationErrorResponse{
			Errors: errors,
		})
	}
	return nil
}

func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", err.Field(), err.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}