package custommiddleware

import "github.com/google/wire"

var Set = wire.NewSet(
	//NewOTELTracer,
	//NewMetrics,
	NewErrorHandling,
	NewTimeout,
	NewAuthentication,
)
