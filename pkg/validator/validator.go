package validator

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	validate *validator.Validate
	// Common regex patterns
	emailRegex    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*]{8,}$`)
	phoneRegex    = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
)

func init() {
	validate = validator.New()

	// Register custom validation tags
	validate.RegisterValidation("email", validateEmail)
	validate.RegisterValidation("password", validatePassword)
	validate.RegisterValidation("phone", validatePhone)

	// Use JSON tags for validation error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validate validates a struct using validator tags
func Validate(i interface{}) error {
	return validate.Struct(i)
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag string) error {
	return validate.Var(field, tag)
}

// Custom validation functions
func validateEmail(fl validator.FieldLevel) bool {
	return emailRegex.MatchString(fl.Field().String())
}

func validatePassword(fl validator.FieldLevel) bool {
	return passwordRegex.MatchString(fl.Field().String())
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors converts validator.ValidationErrors to a slice of ValidationError
func ValidationErrors(err error) []ValidationError {
	var errs []ValidationError

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			errs = append(errs, ValidationError{
				Field:   e.Field(),
				Message: validationErrorMessage(e),
			})
		}
	}

	return errs
}

// validationErrorMessage returns a human-readable error message for a validation error
func validationErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		if e.Type().Kind() == reflect.String {
			return "Must be at least " + e.Param() + " characters long"
		}
		return "Must be at least " + e.Param()
	case "max":
		if e.Type().Kind() == reflect.String {
			return "Must be at most " + e.Param() + " characters long"
		}
		return "Must be at most " + e.Param()
	case "password":
		return "Password must be at least 8 characters long and contain only letters, numbers, and special characters"
	case "phone":
		return "Invalid phone number format"
	default:
		return "Invalid value"
	}
}
