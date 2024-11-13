package custommiddleware

import (
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry/datadog"
)

type DDMetrics struct {
	logger         logger.Logger
	metricsManager *datadog.MetricsManager
}

func NewMetrics(logger logger.Logger, metricsManager *datadog.MetricsManager) *DDMetrics {
	return &DDMetrics{
		logger:         logger,
		metricsManager: metricsManager,
	}
}

func (m *DDMetrics) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := presenter.GetWrapResponseWriter(w)

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			m.metricsManager.RecordHTTPMetrics(
				r.Method,
				r.URL.Path,
				rw.StatusCode,
				duration,
				rw.Length,
			)
		})
	}
}
