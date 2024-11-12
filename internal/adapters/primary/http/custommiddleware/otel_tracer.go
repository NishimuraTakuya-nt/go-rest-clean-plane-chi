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

func OTELTracer() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			tracer := otel.Tracer("http-server")
			resourceName := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			ctx, span := tracer.Start(ctx, resourceName)
			defer span.End()

			start := time.Now()
			log := logger.NewLogger()
			requestID := middleware.GetReqID(ctx)

			// リクエスト開始時のログ
			log.InfoContext(ctx, "Request started")

			// 基本的なリクエスト情報の記録
			span.SetAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.String()),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("http.request_id", requestID),
			)

			// WrapResponseWriter のラッピングを一度だけ行う（後続ではこれを使い回す）
			rw := NewWrapResponseWriter(w)

			// 次のハンドラーの実行
			next.ServeHTTP(rw, r.WithContext(ctx))

			// ルーティングの解決後にルートパターンを取得してspan名を更新
			rctx := chi.RouteContext(ctx)
			if rctx != nil && rctx.RoutePattern() != "" {
				normalizedPattern := fmt.Sprintf("%s %s", r.Method, rctx.RoutePattern())
				span.SetName(normalizedPattern)
				span.SetAttributes(attribute.String("resource.name", normalizedPattern))
			}

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

			log.InfoContext(ctx,
				"Request completed",
				slog.Int("status", rw.StatusCode),
				slog.Int64("bytes", rw.Length),
				slog.String("duration", duration.String()),
			)
		})
	}
}
