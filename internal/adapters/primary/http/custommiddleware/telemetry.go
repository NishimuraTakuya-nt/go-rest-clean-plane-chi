package custommiddleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func Telemetry() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tracer := otel.Tracer("http-server")
			ctx, span := tracer.Start(ctx, fmt.Sprintf("%s %s", r.Method, r.URL.Path))
			defer span.End()

			start := time.Now()
			log := logger.NewLogger()
			requestID := middleware.GetReqID(ctx)

			routeCtx := chi.RouteContext(ctx)
			routePattern := ""
			if routeCtx != nil {
				routePattern = fmt.Sprintf("%s %s", routeCtx.RouteMethod, routeCtx.RoutePattern())
			}

			// リクエスト開始時のログ
			log.InfoContext(r.Context(), "Request started")

			// 基本的なリクエスト情報の記録
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.route_pattern", routePattern),
				attribute.String("http.request_id", requestID),
			)

			// WrapResponseWriter のラッピングを一度だけ行う（後続ではこれを使い回す）
			rw := NewWrapResponseWriter(w)

			// 次のハンドラーの実行
			next.ServeHTTP(rw, r.WithContext(ctx))

			// 処理時間の計算
			duration := time.Since(start)

			// レスポンス情報の記録
			span.SetAttributes(
				attribute.Int("http.status_code", rw.StatusCode),
				attribute.Int64("http.response_size", rw.Length),
				attribute.String("http.duration", duration.String()),
			)

			// エラー情報の記録
			if rw.Err != nil {
				span.RecordError(rw.Err)
				span.SetStatus(codes.Error, rw.Err.Error())
			} else if rw.StatusCode >= 500 {
				// エラーオブジェクトはないが、500番台のステータスコードの場合
				span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", rw.StatusCode))
			}

			log.InfoContext(r.Context(),
				"Request completed",
				slog.Int("status", rw.StatusCode),
				slog.Int64("bytes", rw.Length),
				slog.String("duration", duration.String()),
			)
		})
	}
}
