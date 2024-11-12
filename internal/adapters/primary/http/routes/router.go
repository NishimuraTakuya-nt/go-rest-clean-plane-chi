package routes

import (
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/custommiddleware"
	v1 "github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/routes/v1"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/usecases"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter(
	healthcheckRouter *v1.HealthcheckRouter,
	authRouter *v1.AuthRouter,
	authUsecase usecases.AuthUsecase,
	sampleRouter *v1.SampleRouter,
	telemetryProvider *telemetry.TelemetryProvider,
) (http.Handler, func()) {
	r := chi.NewRouter()
	setupGlobalMiddleware(r, telemetryProvider.Metrics)
	setupSwagger(r)
	setupAPIRoutes(r, healthcheckRouter, authRouter, authUsecase, sampleRouter)
	return r, telemetryProvider.Cleanup
}

func setupGlobalMiddleware(r *chi.Mux, appMetrics *telemetry.AppMetrics) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(custommiddleware.Context())
	r.Use(custommiddleware.OTELTracer())
	//r.Use(custommiddleware.DDTracer()) // fixme choose one of OTELTracer or DDTracer
	r.Use(custommiddleware.Metrics(appMetrics))
	// セキュリティ関連
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   config.Config.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300, // 5 minutes
	}))
	r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
	// APP独自
	r.Use(custommiddleware.ErrorHandler())
	r.Use(custommiddleware.Timeout(config.Config.RequestTimeout))
}

func setupSwagger(r *chi.Mux) {
	if config.Config.Env == "prd" {
		return
	}
	// Swagger 2.0
	r.Get("/swagger/2.0/*", httpSwagger.Handler(httpSwagger.URL("/docs/swagger/swagger.json")))
	// OAS 3.0
	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/docs/swagger/openapi3.json")))
	r.Handle("/docs/swagger/*", http.StripPrefix("/docs/swagger/", http.FileServer(http.Dir("./docs/swagger"))))
}

func setupAPIRoutes(
	r *chi.Mux,
	healthcheckRouter *v1.HealthcheckRouter,
	authRouter *v1.AuthRouter,
	authUsecase usecases.AuthUsecase,
	sampleRouter *v1.SampleRouter,
) {
	r.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {

			r.Group(func(r chi.Router) {
				// 認証不要のパブリックルート
				r.Mount("/healthcheck", healthcheckRouter.Handler)
				r.Mount("/auth", authRouter.Handler)
			})

			r.Group(func(r chi.Router) {
				// 認証必要のプライベートルート
				//r.Use(custommiddleware.Authenticate(authUsecase))
				r.Mount("/samples", sampleRouter.Handler)
			})
		})
	})
}
