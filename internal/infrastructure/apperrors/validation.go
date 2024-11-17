package apperrors

import (
	"fmt"
	"strings"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// FieldError represents a validation error for a specific field
type FieldError struct {
	Field   string
	Value   any
	Message string
}

// ValidationErrors is a collection of FieldErrors
type ValidationErrors []FieldError

// NewValidationErrors creates a new ValidationErrors instance
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{}
}

func (ve *ValidationErrors) Error() string {
	var errMessages []string
	for _, fe := range *ve {
		errMessages = append(errMessages, fmt.Sprintf("%s: %s", fe.Field, fe.Message))
	}
	return strings.Join(errMessages, "; ")
}

// AddError adds a new FieldError to ValidationErrors
func (ve *ValidationErrors) AddError(field string, value any, message string) {
	*ve = append(*ve, FieldError{Field: field, Value: value, Message: message})
}

func (ve *ValidationErrors) AddToSpan(span tracer.Span) {
	if span == nil {
		return
	}
	span.SetTag("error", true)
	span.SetTag("error.type", "validation_error")
	span.SetTag("error.status_code", 400) // バリデーションエラーは通常400
	span.SetTag("error.message", ve.Error())
	span.SetTag("validation.error_count", len(*ve))
}
