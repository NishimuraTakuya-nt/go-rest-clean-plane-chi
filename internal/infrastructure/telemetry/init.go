package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

type OpenTelemetryProvider struct {
	Metrics *AppMetrics
	Cleanup func()
}

// InitTelemetry テレメトリーの初期化とTelemetryProviderの作成
func InitTelemetry(cfg *config.AppConfig, logger logger.Logger) (*OpenTelemetryProvider, error) {
	ctx := context.Background()

	// リソース情報の設定
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("go-rest-clean-plane-chi"),
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

	// トレーサーの初期化
	tp, err := initTracer(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	// メトリクスの初期化
	mp, metrics, err := initMetrics(ctx, res, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize metrics: %w", err)
	}

	// クリーンアップ関数
	cleanup := func() {
		logger.Info("Shutting down telemetry provider & meter provider")
		ctx := context.Background()
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down tracer provider", "error", err.Error())
			return
		}
		if err := mp.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down meter provider", "error", err.Error())
			return
		}
		logger.Info("Telemetry provider & meter provider shut down successfully")
	}

	provider := &OpenTelemetryProvider{
		Metrics: metrics,
		Cleanup: cleanup,
	}

	return provider, nil
}

// initTracer トレーサーの初期化
func initTracer(ctx context.Context, res *resource.Resource, cfg *config.AppConfig) (*sdktrace.TracerProvider, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.DDAgentHost, cfg.DDAgentPort)
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithTimeout(5*time.Second),
	)

	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter,
			sdktrace.WithBatchTimeout(5*time.Second),
			sdktrace.WithMaxExportBatchSize(512),
			sdktrace.WithMaxQueueSize(2048),
		),
		sdktrace.WithResource(res),
		// サンプリング設定の詳細化 // fixme 一旦１100%にしてDatadogで確認しやすくしおく
		//sdktrace.WithSampler(sdktrace.ParentBased(
		//	sdktrace.TraceIDRatioBased(1.0), // 10%のサンプリング
		//)),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// グローバルトレーサープロバイダーの設定
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

// initMetrics メトリクスの初期化
func initMetrics(ctx context.Context, res *resource.Resource, cfg *config.AppConfig) (*sdkmetric.MeterProvider, *AppMetrics, error) {
	endpoint := fmt.Sprintf("%s:%s", cfg.DDAgentHost, cfg.DDAgentPort)
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(endpoint),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metric exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exporter,
				sdkmetric.WithInterval(30*time.Second),
			),
		),
	)
	otel.SetMeterProvider(meterProvider)
	meter := meterProvider.Meter("app-metrics")
	metrics, err := newAppMetrics(meter)
	return meterProvider, metrics, err
}

func getHostName() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
