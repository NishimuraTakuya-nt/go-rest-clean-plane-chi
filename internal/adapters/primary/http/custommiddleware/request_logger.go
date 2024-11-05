package custommiddleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

// RequestLogger logs details of each HTTP request
func RequestLogger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log := logger.NewLogger()

			// リクエスト開始時のログ
			log.InfoContext(r.Context(), "Request started")

			// ResponseWriter のラッピングを一度だけ行う（後続のMiddlewareではこれを使い回す）
			rw := NewResponseWriter(w)

			defer func() {
				log.InfoContext(r.Context(),
					"Request completed",
					slog.Int("status", rw.StatusCode),
					slog.Int64("bytes", rw.Length),
					slog.String("duration", time.Since(start).String()),
				)
			}()

			// 次のハンドラを呼び出し
			next.ServeHTTP(rw, r)
		})
	}
}
