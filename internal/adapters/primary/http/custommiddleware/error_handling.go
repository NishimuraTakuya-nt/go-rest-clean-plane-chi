package custommiddleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/dto/response"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/adapters/primary/http/presenter"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/apperrors"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/go-chi/chi/v5/middleware"
)

type ErrorHandling struct {
	logger     logger.Logger
	JSONWriter *presenter.JSONWriter
}

func NewErrorHandling(
	logger logger.Logger,
	JSONWriter *presenter.JSONWriter,
) *ErrorHandling {
	return &ErrorHandling{
		logger:     logger,
		JSONWriter: JSONWriter,
	}
}

func (h *ErrorHandling) Handle() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := presenter.GetWrapResponseWriter(w)

			defer func() {
				if re := recover(); re != nil {
					var panicErr error
					switch err := re.(type) {
					case string:
						panicErr = errors.New(err)
					case error:
						panicErr = err
					default:
						panicErr = fmt.Errorf("unknown panic: %v", err)
					}

					// スタックトレースを取得
					stack := debug.Stack()

					// 詳細なログを記録
					h.logger.ErrorContext(r.Context(),
						"Panic occurred",
						"error", panicErr,
						"stack", string(stack),
					)

					// クライアントへのレスポンス用エラー
					clientErr := apperrors.NewInternalError("Unexpected error occurred", panicErr)

					// エラーハンドリング
					h.handleError(r.Context(), rw, clientErr)
				}
			}()

			next.ServeHTTP(rw, r)

			if rw.Err != nil {
				h.handleError(r.Context(), rw, rw.Err)
			}
		})
	}
}

func (h *ErrorHandling) handleError(ctx context.Context, rw *presenter.WrapResponseWriter, err error) {
	var res response.ErrorResponse
	var statusCode int
	requestID := middleware.GetReqID(ctx)

	switch e := err.(type) {
	case *apperrors.ValidationErrors:
		statusCode = http.StatusBadRequest
		details := make([]map[string]any, 0, len(*e))
		for _, fe := range *e {
			details = append(details, map[string]any{
				"field":   fe.Field,
				"value":   fe.Value,
				"message": fe.Message,
			})
		}
		res = response.ErrorResponse{
			StatusCode: statusCode,
			Type:       string(apperrors.ErrorTypeBadRequest),
			RequestID:  requestID,
			Message:    "Validation error",
			Details:    details,
		}

	case *apperrors.AppError:
		statusCode = e.StatusCode
		res = response.ErrorResponse{
			StatusCode: statusCode,
			Type:       string(e.Type),
			RequestID:  requestID,
			Message:    e.Message,
		}

	default:
		statusCode = http.StatusInternalServerError
		res = response.ErrorResponse{
			StatusCode: statusCode,
			Type:       string(apperrors.ErrorTypeInternal),
			RequestID:  requestID,
			Message:    "Internal server error",
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	h.JSONWriter.Write(ctx, rw, res)
}
