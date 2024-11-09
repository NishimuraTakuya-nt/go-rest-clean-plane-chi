package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func InitTelemetry(ctx context.Context) (func(context.Context) error, error) {
	log := logger.NewLogger()
	cfg := config.Config

	// リソース情報の設定
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.Version),
			semconv.DeploymentEnvironment(cfg.Env),
			semconv.HostName(getHostName()),
		),
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithContainer(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// エクスポーターの設定
	endpoint := fmt.Sprintf("%s:%s", cfg.DDAgentHost, cfg.DDAgentPort)
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		// タイムアウトの設定
		otlptracegrpc.WithTimeout(5*time.Second),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// トレーサープロバイダーの設定
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
			// エラー時の再試行設定
			sdktrace.WithMaxQueueSize(2048),
		),
		sdktrace.WithResource(res),
		// サンプリング設定の詳細化
		sdktrace.WithSampler(sdktrace.ParentBased(
			sdktrace.TraceIDRatioBased(0.1), // 10%のサンプリング
		)),
	)

	// グローバルトレーサープロバイダーの設定
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		log.Info("Shutting down telemetry provider")
		if err := tp.Shutdown(ctx); err != nil {
			log.Error("Error shutting down tracer provider", "error", err.Error())
			return fmt.Errorf("failed to shutdown tracer provider: %w", err)
		}
		log.Info("Telemetry provider shut down successfully")
		return nil
	}, nil
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
