//go:build wireinject
// +build wireinject

package main

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/routes"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/routes/v1"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/secondary/piyographql"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/usecases"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/auth"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry"
	"github.com/google/wire"
)

func InitializeAPI() (http.Handler, func(), error) {
	wire.Build(
		logger.Set,
		piyographql.Set,
		auth.Set,
		usecases.Set,
		handlers.Set,
		v1.Set,
		telemetry.Set,
		routes.Set,
	)
	return nil, nil, nil
}
