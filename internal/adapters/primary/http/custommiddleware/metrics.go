package custommiddleware

import (
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry"
)

func Metrics(metrics *telemetry.AppMetrics) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()
			rw := GetWrapResponseWriter(w)

			next.ServeHTTP(rw, r)

			duration := time.Since(startTime)
			metrics.RecordHTTPRequest(
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
