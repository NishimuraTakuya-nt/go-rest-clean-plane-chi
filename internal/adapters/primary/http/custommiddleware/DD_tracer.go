package custommiddleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/apperrors"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type DDTracer struct {
	logger logger.Logger
}

func NewDDTracer(logger logger.Logger) *DDTracer {
	return &DDTracer{
		logger: logger,
	}
}

func (t *DDTracer) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()

			opts := []tracer.StartSpanOption{
				tracer.ResourceName(fmt.Sprintf("%s %s", r.Method, r.URL.Path)),
				tracer.SpanType("web"),
				tracer.Tag("http.method", r.Method),
				tracer.Tag("http.url", r.URL.String()),
				tracer.Tag("http.user_agent", r.UserAgent()),
				tracer.Tag("request_id", middleware.GetReqID(ctx)),
			}

			span, ctx := tracer.StartSpanFromContext(ctx, "http.request", opts...)
			defer span.Finish()

			// リクエスト開始時のログ
			t.logger.InfoContext(ctx, "Request started")

			// WrapResponseWriter のラッピングを一度だけ行う（後続ではこれを使い回す）
			rw := presenter.NewWrapResponseWriter(w)

			next.ServeHTTP(rw, r.WithContext(ctx))

			duration := time.Since(start)

			// ルーティングパターンの取得と設定
			rctx := chi.RouteContext(ctx)
			if rctx != nil && rctx.RoutePattern() != "" {
				normalizedPattern := fmt.Sprintf("%s %s", r.Method, rctx.RoutePattern())
				// スパン名の更新
				span.SetOperationName(normalizedPattern)
				// リソース名の更新
				span.SetTag("resource.name", normalizedPattern)
			}

			// レスポンス情報の記録
			span.SetTag("http.status_code", rw.StatusCode)
			span.SetTag("http.response_size", rw.Length)
			span.SetTag("http.duration", duration.String())

			if rw.Err != nil {
				var errorTracer apperrors.ErrorTracer
				if errors.As(rw.Err, &errorTracer) {
					// Spanにエラー情報を追加
					errorTracer.AddToSpan(span)

				}
			} else if rw.StatusCode >= 500 {
				span.SetTag("error", true)
				span.SetTag("error.message", fmt.Sprintf("HTTP %d", rw.StatusCode))
				span.SetTag("error.type", "http_error")
			}

			t.logger.InfoContext(ctx,
				"Request completed",
				slog.Int("status", rw.StatusCode),
				slog.Int64("bytes", rw.Length),
				slog.String("duration", duration.String()),
			)
		})
	}
}
