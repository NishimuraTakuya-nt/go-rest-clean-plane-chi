package custommiddleware

import (
	"context"
	"net/http"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/apperrors"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

type Timeout struct {
	logger logger.Logger
	cfg    *config.AppConfig
}

func NewTimeout(
	logger logger.Logger,
	cfg *config.AppConfig,
) *Timeout {
	return &Timeout{
		logger: logger,
		cfg:    cfg,
	}
}

func (h Timeout) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), h.cfg.RequestTimeout)
			defer cancel()
			rw := presenter.GetWrapResponseWriter(w)

			done := make(chan bool)
			go func() {
				next.ServeHTTP(rw, r.WithContext(ctx))
				done <- true
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				h.logger.ErrorContext(r.Context(), "Request timed out")
				rw.WriteError(apperrors.NewTimeoutError("Request timed out", ctx.Err()))
			}
		})
	}
}
