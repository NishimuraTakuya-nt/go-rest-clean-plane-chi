package presenter

import "github.com/google/wire"

var Set = wire.NewSet(
	NewWrapResponseWriter,
	NewJSONWriter,
)
