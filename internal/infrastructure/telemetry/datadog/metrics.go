package datadog

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
)

type MetricsManager struct {
	logger         logger.Logger
	client         *statsd.Client
	prefix         string
	tags           []string
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.Mutex
	metricsBuffer  chan metricEvent
	bufferSize     int
	flushInterval  time.Duration
	shutdownSignal chan struct{}
}

// メトリクスイベントを表す構造体
type metricEvent struct {
	metricType string
	name       string
	value      float64
	tags       []string
	rate       float64
}

const (
	defaultBufferSize    = 1000
	defaultFlushInterval = 1 * time.Second
	shutdownTimeout      = 5 * time.Second
)

func NewMetricsManager(cfg *config.AppConfig, logger logger.Logger) (*MetricsManager, error) {
	if !cfg.DDEnabled {
		logger.Info("Datadog metrics is disabled")
		return &MetricsManager{}, nil
	}

	client, err := statsd.New(
		fmt.Sprintf("%s:%s", cfg.DDAgentHost, cfg.DDAgentMetricsPort),
		statsd.WithBufferPoolSize(defaultBufferSize),
		statsd.WithMaxMessagesPerPayload(defaultBufferSize/2),
		statsd.WithNamespace("app."),
		statsd.WithTags([]string{
			fmt.Sprintf("env:%s", cfg.Env),
			fmt.Sprintf("version:%s", cfg.Version),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create StatsD client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	mm := &MetricsManager{
		logger:         logger,
		client:         client,
		prefix:         "app.",
		ctx:            ctx,
		cancel:         cancel,
		metricsBuffer:  make(chan metricEvent, defaultBufferSize),
		bufferSize:     defaultBufferSize,
		flushInterval:  defaultFlushInterval,
		shutdownSignal: make(chan struct{}),
	}

	// メトリクス処理用のワーカーを起動
	go mm.processMetrics()

	return mm, nil
}

func (m *MetricsManager) Start() {
	if m.client == nil {
		return
	}
	go m.collectSystemMetrics()
}

func (m *MetricsManager) Stop() {
	if m.client == nil {
		return
	}

	m.logger.Info("Stopping metrics manager")
	m.cancel()

	// bufferの処理を待つ
	select {
	case <-m.shutdownSignal:
		m.logger.Info("Metrics manager shutdownSignal received")
	case <-time.After(shutdownTimeout):
		m.logger.Warn("Metrics manager shutdown timed out")
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if err := m.client.Close(); err != nil {
		m.logger.Error("Error closing StatsD client", "error", err)
	} else {
		m.logger.Info("Metrics manager shutdown completed")
	}
}

func (m *MetricsManager) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// メモリ統計用の構造体を再利用
	var mem runtime.MemStats

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("Stopping system metrics collection")
			return
		case <-ticker.C:
			runtime.ReadMemStats(&mem)

			systemMetrics := []metricEvent{
				{
					metricType: "gauge",
					name:       "system.goroutines",
					value:      float64(runtime.NumGoroutine()),
					rate:       1.0,
				},
				{
					metricType: "gauge",
					name:       "system.memory.allocated",
					value:      float64(mem.Alloc),
					tags:       []string{"type:heap"},
					rate:       1.0,
				},
				{
					metricType: "gauge",
					name:       "system.memory.total",
					value:      float64(mem.TotalAlloc),
					tags:       []string{"type:total"},
					rate:       1.0,
				},
			}

			for _, metric := range systemMetrics {
				select {
				case m.metricsBuffer <- metric:
				default:
					m.logger.Warn("Metrics buffer is full, dropping system metric",
						"type", metric.metricType,
						"name", metric.name)
				}
			}
		}
	}
}

func (m *MetricsManager) RecordHTTPMetrics(method, path string, statusCode int, duration time.Duration, responseSize int64) {
	if m.client == nil {
		return
	}

	tags := []string{
		fmt.Sprintf("method:%s", method),
		fmt.Sprintf("path:%s", path),
		fmt.Sprintf("status:%d", statusCode),
	}

	metrics := []metricEvent{
		{
			metricType: "count",
			name:       "http.request.count",
			value:      1,
			tags:       tags,
			rate:       1.0,
		},
		{
			metricType: "gauge",
			name:       "http.request.duration",
			value:      float64(duration.Milliseconds()),
			tags:       tags,
			rate:       1.0,
		},
		{
			metricType: "histogram",
			name:       "http.response.size",
			value:      float64(responseSize),
			tags:       tags,
			rate:       1.0,
		},
	}

	// エラーメトリクスを条件付きで追加
	if statusCode >= 400 {
		metrics = append(metrics, metricEvent{
			metricType: "count",
			name:       "http.request.errors",
			value:      1,
			tags:       tags,
			rate:       1.0,
		})
	}

	// メトリクスをバッファに追加
	for _, metric := range metrics {
		select {
		case m.metricsBuffer <- metric:
		default:
			m.logger.Warn("Metrics buffer is full, dropping metric",
				"type", metric.metricType,
				"name", metric.name)
		}
	}
}

// flush worker
func (m *MetricsManager) processMetrics() {
	ticker := time.NewTicker(m.flushInterval)
	defer ticker.Stop()
	defer close(m.shutdownSignal)

	buffer := make([]metricEvent, 0, m.bufferSize)

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("Processing final metrics before shutdown")
			if len(buffer) > 0 {
				// 残りのメトリクスを処理
				m.flushMetrics(buffer)
			}
			return

		case event := <-m.metricsBuffer:
			// イベントをbuffer変数に追加
			buffer = append(buffer, event)
			if len(buffer) >= m.bufferSize {
				// バッファが一杯になったらメトリクスを送信
				m.flushMetrics(buffer)
				buffer = buffer[:0]
			}

		case <-ticker.C:
			// インターバルごとにbuffer変数をメトリクスを送信
			if len(buffer) > 0 {
				m.flushMetrics(buffer)
				buffer = buffer[:0]
			}
		}
	}
}

func (m *MetricsManager) flushMetrics(buffer []metricEvent) {
	if m.client == nil {
		return
	}

	for _, event := range buffer {
		var err error
		switch event.metricType {
		case "gauge":
			err = m.client.Gauge(event.name, event.value, event.tags, event.rate)
		case "count":
			err = m.client.Count(event.name, int64(event.value), event.tags, event.rate)
		case "histogram":
			err = m.client.Histogram(event.name, event.value, event.tags, event.rate)
		default:
			m.logger.Error("Unknown metric type",
				"type", event.metricType,
				"name", event.name)
			continue
		}

		if err != nil {
			m.logger.Error("Failed to send metric",
				"type", event.metricType,
				"name", event.name,
				"error", err)
		}
	}
}
