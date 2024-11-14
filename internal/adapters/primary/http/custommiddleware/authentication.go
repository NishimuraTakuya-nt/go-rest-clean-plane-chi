package custommiddleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/apperrors"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/usecases"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

var UserKey = struct{}{}

var excludedPaths = []string{
	"/api/v1/auth/login",
	"/api/v1/healthcheck",
	"/swagger/",
	"/docs/swagger/",
}

type Authentication struct {
	logger      logger.Logger
	JSONWriter  *presenter.JSONWriter
	authUsecase usecases.AuthUsecase
}

func NewAuthentication(
	logger logger.Logger,
	JSONWriter *presenter.JSONWriter,
	authUsecase usecases.AuthUsecase,
) *Authentication {
	return &Authentication{
		logger:      logger,
		JSONWriter:  JSONWriter,
		authUsecase: authUsecase,
	}
}

func (h *Authentication) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := presenter.GetWrapResponseWriter(w)

			// 除外パスのチェック
			for _, path := range excludedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(rw, r)
					return
				}
			}

			// 認証方法は適宜変更する
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				h.logger.ErrorContext(r.Context(), "Missing authorization header")
				rw.WriteError(apperrors.NewUnauthorizedError("Missing authorization header", nil))
				return
			}

			if len(authHeader) <= 7 || authHeader[:7] != "Bearer " {
				h.logger.ErrorContext(r.Context(), "Invalid token format", "header", authHeader)
				rw.WriteError(apperrors.NewUnauthorizedError("Invalid token format", nil))
				return
			}

			tokenString := authHeader[7:]
			user, err := h.authUsecase.Authenticate(r.Context(), tokenString)
			if err != nil {
				h.logger.ErrorContext(r.Context(), "Token validation failed", "error", err)
				rw.WriteError(apperrors.NewUnauthorizedError("Invalid or expired token", nil))
				return
			}

			// nolint:staticcheck
			ctx := context.WithValue(r.Context(), UserKey, user)
			h.logger.InfoContext(r.Context(), "User authenticated")
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
