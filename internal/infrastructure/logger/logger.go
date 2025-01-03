package logger

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/common/contextkeys"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/go-chi/chi/v5/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Logger interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)
	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
	With(args ...any) Logger
}

type customLogger struct {
	logger *slog.Logger
}

func NewLogger(cfg *config.AppConfig) Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     getLevelFromEnv(cfg.LogLevel),
		ReplaceAttr: func(_ []string, att slog.Attr) slog.Attr {
			if att.Key == slog.SourceKey {
				const skip = 7
				_, file, line, ok := runtime.Caller(skip)
				if !ok {
					return att
				}
				v := fmt.Sprintf("%s:%d", filepath.Base(file), line)
				att.Value = slog.StringValue(v)
			}
			return att
		},
	})

	return &customLogger{
		logger: slog.New(handler),
	}
}

// addStandardFields 標準フィールドを追加
func (l *customLogger) addStandardFields(ctx context.Context, args []any) []any {
	if ctx != nil {
		if requestID := middleware.GetReqID(ctx); requestID != "" {
			args = append(args, "request_id", requestID)
		}
		if u, ok := ctx.Value(contextkeys.UserIDKey).(string); ok {
			args = append(args, "user_id", u)
		}
		if r, ok := ctx.Value(contextkeys.HTTPRequestKey).(*http.Request); ok {
			args = append(args,
				"endpoint", r.URL.Path,
				"method", r.Method,
				"remote_ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"referer", r.Referer(),
			)
		}
		// トレース情報の追加
		//if span := trace.SpanFromContext(ctx); span.IsRecording() { // case: opentelemetry
		//	spanContext := span.SpanContext()
		//	if spanContext.HasTraceID() {
		//		args = append(args,
		//			"trace_id", spanContext.TraceID().String(),
		//			"span_id", spanContext.SpanID().String(),
		//		)
		//	}
		//}
		if span, ok := tracer.SpanFromContext(ctx); ok {
			spanCtx := span.Context()
			args = append(args,
				"dd.trace_id", fmt.Sprintf("%d", spanCtx.TraceID()),
				"dd.span_id", fmt.Sprintf("%d", spanCtx.SpanID()),
			)
		}
	}
	return args
}

func (l *customLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *customLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	args = l.addStandardFields(ctx, args)
	l.logger.DebugContext(ctx, msg, args...)
}

func (l *customLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *customLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	args = l.addStandardFields(ctx, args)
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *customLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *customLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	args = l.addStandardFields(ctx, args)
	l.logger.WarnContext(ctx, msg, args...)
}

func (l *customLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *customLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	args = l.addStandardFields(ctx, args)
	l.logger.ErrorContext(ctx, msg, args...)
}

// With 追加のフィールドを持つ新しいロガーを作成
func (l *customLogger) With(args ...any) Logger {
	return &customLogger{
		logger: l.logger.With(args...),
	}
}

func getLevelFromEnv(logLevel string) slog.Level {
	switch logLevel {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
