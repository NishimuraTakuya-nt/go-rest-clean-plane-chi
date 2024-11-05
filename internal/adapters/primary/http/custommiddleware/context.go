package custommiddleware

import (
	"context"
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/common/contextkeys"
	"github.com/go-chi/chi/v5/middleware"
)

func Context() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// リクエストIDをheaderに追加
			requestID := middleware.GetReqID(ctx)
			w.Header().Set("X-Request-ID", requestID)

			// リクエスト情報をコンテキストに追加
			ctx = context.WithValue(ctx, contextkeys.HTTPRequestKey, r)

			// 次のハンドラを呼び出し
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
