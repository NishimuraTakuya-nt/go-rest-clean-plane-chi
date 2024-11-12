package telemetry

import "github.com/google/wire"

var Set = wire.NewSet(
	InitTelemetry,
)
