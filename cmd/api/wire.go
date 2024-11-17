//go:build wireinject
// +build wireinject

package main

import (
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/custommiddleware"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/handlers"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/routes"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/routes/v1"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/secondary/piyographql"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/services"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/usecases"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry/datadog"
	"github.com/google/wire"
)

func InitializeRouter(cfg *config.AppConfig, logger logger.Logger, metricsManager *datadog.MetricsManager) (*routes.Router, error) {
	wire.Build(
		presenter.Set,
		custommiddleware.Set,
		piyographql.Set,
		services.Set,
		usecases.Set,
		handlers.Set,
		v1.Set,
		//telemetry.Set,
		routes.Set,
	)
	return nil, nil
}
