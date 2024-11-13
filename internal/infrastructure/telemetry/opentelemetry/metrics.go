package opentelemetry

import (
	"context"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// AppMetrics アプリケーションのメトリクスを管理
type AppMetrics struct {
	// システムメトリクス
	goroutineCount metric.Int64UpDownCounter
	memoryUsage    metric.Int64UpDownCounter

	// HTTPメトリクス
	requestCount    metric.Int64Counter
	requestDuration metric.Float64Histogram
	responseSize    metric.Int64Histogram

	// カスタムメトリクス
	businessOpCount metric.Int64Counter
}

// NewAppMetrics メトリクスの初期化
func newAppMetrics(meter metric.Meter) (*AppMetrics, error) {
	goroutineCount, err := meter.Int64UpDownCounter(
		"app.goroutines",
		metric.WithDescription("Number of goroutines"),
	)
	if err != nil {
		return nil, err
	}

	memoryUsage, err := meter.Int64UpDownCounter(
		"app.memory.usage",
		metric.WithDescription("Memory usage in bytes"),
	)
	if err != nil {
		return nil, err
	}

	requestCount, err := meter.Int64Counter(
		"app.http.request.count",
		metric.WithDescription("Total number of HTTP requests"),
	)
	if err != nil {
		return nil, err
	}

	requestDuration, err := meter.Float64Histogram(
		"app.http.request.duration",
		metric.WithDescription("HTTP request duration"),
		metric.WithUnit("ms"),
	)
	if err != nil {
		return nil, err
	}

	responseSize, err := meter.Int64Histogram(
		"app.http.response.size",
		metric.WithDescription("HTTPレスポンスサイズ"),
		metric.WithUnit("bytes"),
	)
	if err != nil {
		return nil, err
	}

	businessOpCount, err := meter.Int64Counter(
		"app.business.operation.count",
		metric.WithDescription("ビジネスオペレーション実行数"),
	)
	if err != nil {
		return nil, err
	}

	metrics := &AppMetrics{
		goroutineCount:  goroutineCount,
		memoryUsage:     memoryUsage,
		requestCount:    requestCount,
		requestDuration: requestDuration,
		responseSize:    responseSize,
		businessOpCount: businessOpCount,
	}

	// システムメトリクスの収集開始
	go metrics.collectSystemMetrics(context.Background())

	return metrics, nil
}

// collectSystemMetrics システムメトリクスの定期収集
func (am *AppMetrics) collectSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			am.memoryUsage.Add(ctx, int64(m.Alloc),
				metric.WithAttributes(
					attribute.String("type", "heap"),
					attribute.String("unit", "bytes"),
				))

			am.goroutineCount.Add(ctx, int64(runtime.NumGoroutine()))

		case <-ctx.Done():
			return
		}
	}
}

// RecordHTTPRequest HTTPリクエストのメトリクスを記録
func (am *AppMetrics) RecordHTTPRequest(ctx context.Context, method, path string, duration time.Duration, statusCode int, responseSize int64) {
	attrs := []attribute.KeyValue{
		attribute.String("method", method),
		attribute.String("path", path),
		attribute.Int("status_code", statusCode),
	}

	am.requestCount.Add(ctx, 1, metric.WithAttributes(attrs...))
	am.requestDuration.Record(ctx, float64(duration.Milliseconds()), metric.WithAttributes(attrs...))
	am.responseSize.Record(ctx, responseSize, metric.WithAttributes(attrs...))
}

// RecordBusinessOperation ビジネスオペレーションのメトリクスを記録
func (am *AppMetrics) RecordBusinessOperation(ctx context.Context, operationType string, success bool) {
	attrs := []attribute.KeyValue{
		attribute.String("type", operationType),
		attribute.Bool("success", success),
	}

	am.businessOpCount.Add(ctx, 1, metric.WithAttributes(attrs...))
}
