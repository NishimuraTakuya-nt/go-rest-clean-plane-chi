package custommiddleware

import (
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry/opentelemetry"
)

type OTELMetrics struct {
	metrics *opentelemetry.AppMetrics
}

func NewOTELMetrics(metrics *opentelemetry.AppMetrics) *OTELMetrics {
	return &OTELMetrics{metrics: metrics}
}

func (h OTELMetrics) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			rw := presenter.GetWrapResponseWriter(w)

			next.ServeHTTP(rw, r)

			duration := time.Since(startTime)
			h.metrics.RecordHTTPRequest(
				r.Context(),
				r.Method,
				r.URL.Path,
				duration,
				rw.StatusCode,
				rw.Length,
			)
		})
	}
}
