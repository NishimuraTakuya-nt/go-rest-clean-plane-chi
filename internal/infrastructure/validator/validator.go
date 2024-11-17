package validator

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/apperrors"
	"github.com/go-playground/validator/v10"
)

var (
	once     sync.Once
	validate *validator.Validate
)

// initValidator initializes the validator instance with custom validations
func initValidator() {
	validate = validator.New()

	// Register custom validation for sample ID
	_ = validate.RegisterValidation("sampleId", sampleID)

	// Add more custom validations here as needed
}

// GetValidator returns a singleton instance of the validator
func GetValidator() *validator.Validate {
	once.Do(initValidator)
	return validate
}

// Validate validates a struct based on the validator
func Validate(s interface{}) *apperrors.ValidationErrors {
	err := GetValidator().Struct(s)
	if err == nil {
		return nil
	}

	errors := apperrors.NewValidationErrors()
	for _, err := range err.(validator.ValidationErrors) {
		errors.AddError(err.Namespace(), err.Value(), getErrorMsg(err))
	}
	return errors
}

// ValidateVar validates a single variable
func ValidateVar(field interface{}, tag, namespace string) *apperrors.ValidationErrors {
	err := GetValidator().Var(field, tag)
	if err == nil {
		return nil
	}

	errors := apperrors.NewValidationErrors()
	for _, err := range err.(validator.ValidationErrors) {
		errors.AddError(namespace, err.Value(), getErrorMsg(err))
	}
	return errors
}

// sampleID is a custom validation function for sample IDs
func sampleID(fl validator.FieldLevel) bool {
	// sampleID must be alphanumeric and between 3 and 20 characters
	return regexp.MustCompile(`^[a-zA-Z0-9]{3,20}$`).MatchString(fl.Field().String())
}

func getErrorMsg(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return fmt.Sprintf("Minimum length is %s", err.Param())
	case "max":
		return fmt.Sprintf("Maximum length is %s", err.Param())
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", err.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", err.Param())
	default:
		return "Invalid value"
	}
}
