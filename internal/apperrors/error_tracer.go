package apperrors

import "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

type ErrorTracer interface {
	error
	AddToSpan(span tracer.Span)
}
