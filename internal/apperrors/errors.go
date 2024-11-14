package apperrors

import (
	"net/http"
	"runtime"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type ErrorType string

const (
	ErrorTypeBadRequest         ErrorType = "BAD_REQUEST"
	ErrorTypeUnauthorized       ErrorType = "UNAUTHORIZED"
	ErrorTypeForbidden          ErrorType = "FORBIDDEN"
	ErrorTypeNotFound           ErrorType = "NOT_FOUND"
	ErrorTypeConflict           ErrorType = "CONFLICT"
	ErrorTypeRateLimit          ErrorType = "RATE_LIMIT"
	ErrorTypeInternal           ErrorType = "INTERNAL_ERROR"
	ErrorTypeExternalService    ErrorType = "EXTERNAL_SERVICE_ERROR"
	ErrorTypeServiceUnavailable ErrorType = "SERVICE_UNAVAILABLE"
	ErrorTypeTimeout            ErrorType = "TIMEOUT"
)

// AppError はアプリケーション固有のエラーを表します。
type AppError struct {
	Type       ErrorType
	RawError   error
	StatusCode int
	Message    string
	File       string
	Line       int
	Function   string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(errType ErrorType, rawErr error, statusCode int, message string) *AppError {
	appErr := &AppError{
		Type:       errType,
		RawError:   rawErr,
		StatusCode: statusCode,
		Message:    message,
	}
	appErr.captureStack(2) // 呼び出し元の情報を取得
	return appErr
}

func (e *AppError) captureStack(skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	if !ok {
		return
	}

	fn := runtime.FuncForPC(pc)
	if fn != nil {
		e.Function = fn.Name()
	}
	e.File = file
	e.Line = line
}

func (e *AppError) AddToSpan(span tracer.Span) {
	if span == nil {
		return
	}
	span.SetTag("error", true)
	span.SetTag("error.type", string(e.Type))
	span.SetTag("error.message", e.Message)
	span.SetTag("error.status_code", e.StatusCode)
	span.SetTag("error.file", e.File)
	span.SetTag("error.line", e.Line)
	span.SetTag("error.function", e.Function)

	if e.RawError != nil {
		span.SetTag("error.raw", e.RawError.Error())
	}
}

// NewBadRequestError 400 Bad Request
func NewBadRequestError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeBadRequest, rawErr, http.StatusBadRequest, message)
}

// NewUnauthorizedError 401 Unauthorized
func NewUnauthorizedError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeUnauthorized, rawErr, http.StatusUnauthorized, message)
}

// NewForbiddenError 403 Forbidden
func NewForbiddenError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeForbidden, rawErr, http.StatusForbidden, message)
}

// NewNotFoundError 404 Not Found
func NewNotFoundError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeNotFound, rawErr, http.StatusNotFound, message)
}

// NewConflictError 409 Conflict
func NewConflictError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeConflict, rawErr, http.StatusConflict, message)
}

// NewRateLimitError 429 Too Many Requests
func NewRateLimitError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeRateLimit, rawErr, http.StatusTooManyRequests, message)
}

// NewInternalError 500 Internal Server Error
func NewInternalError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeInternal, rawErr, http.StatusInternalServerError, message)
}

// NewExternalServiceError 502 Bad Gateway
func NewExternalServiceError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeExternalService, rawErr, http.StatusBadGateway, message)
}

// NewServiceUnavailableError 503 Service Unavailable
func NewServiceUnavailableError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeServiceUnavailable, rawErr, http.StatusServiceUnavailable, message)
}

// NewTimeoutError 504 Gateway Timeout
func NewTimeoutError(message string, rawErr error) *AppError {
	return NewAppError(ErrorTypeTimeout, rawErr, http.StatusGatewayTimeout, message)
}
