package custommiddleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
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

			// リクエスト開始時のログ
			t.logger.InfoContext(ctx, "Request started")

			// スパンの作成
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

			// WrapResponseWriter のラッピングを一度だけ行う（後続ではこれを使い回す）
			rw := presenter.NewWrapResponseWriter(w)

			// 次のハンドラーの実行
			next.ServeHTTP(rw, r.WithContext(ctx))

			// 処理時間の計算
			duration := time.Since(start)

			// ルーティングパターンの取得と設定
			rctx := chi.RouteContext(ctx)
			if rctx != nil && rctx.RoutePattern() != "" {
				normalizedPattern := fmt.Sprintf("%s %s", r.Method, rctx.RoutePattern())
				// スパン名の更新
				span.SetOperationName(normalizedPattern)
				// リソース名の設定
				span.SetTag("resource.name", normalizedPattern)
			}

			// レスポンス情報の記録
			span.SetTag("http.status_code", rw.StatusCode)
			span.SetTag("http.response_size", rw.Length)
			span.SetTag("http.duration", duration.String())

			// エラーハンドリング
			if rw.Err != nil {
				span.SetTag("error", true)
				span.SetTag("error.message", rw.Err.Error())
				span.SetTag("error.type", fmt.Sprintf("%T", rw.Err))
				// スタックトレースが必要な場合
				if stackTracer, ok := rw.Err.(interface{ StackTrace() []byte }); ok {
					span.SetTag("error.stack", string(stackTracer.StackTrace()))
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
