package custommiddleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/go-chi/chi/v5/middleware"
	chitrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-chi/chi.v5"
)

func DDTracer() Middleware {
	// Datadogのミドルウェア
	ddMiddleware := chitrace.Middleware(
		chitrace.WithServiceName("go-rest-clean-plane-chi"),
		chitrace.WithAnalytics(true),
	)

	return func(next http.Handler) http.Handler {
		// Datadogミドルウェアでラップ
		wrappedHandler := ddMiddleware(next)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			log := logger.NewLogger()
			ctx := r.Context()
			requestID := middleware.GetReqID(ctx)

			// リクエスト開始時のログ
			log.InfoContext(ctx, "Request started")

			rw := NewWrapResponseWriter(w)

			// Datadogでラップされたハンドラーを実行
			wrappedHandler.ServeHTTP(rw, r.WithContext(ctx))

			// 処理時間の計算
			duration := time.Since(start)

			log.InfoContext(ctx,
				"Request completed",
				slog.String("request_id", requestID),
				slog.Int("status", rw.StatusCode),
				slog.Int64("bytes", rw.Length),
				slog.String("duration", duration.String()),
			)
		})
	}
}
