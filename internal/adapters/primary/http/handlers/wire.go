package handlers

import "github.com/google/wire"

var Set = wire.NewSet(
	NewHealthcheckHandler,
	NewAuthHandler,
	NewSampleHandler,
)
