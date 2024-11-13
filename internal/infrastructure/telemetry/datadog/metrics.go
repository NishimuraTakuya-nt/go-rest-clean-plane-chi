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
	logger logger.Logger
	client *statsd.Client
	prefix string
	tags   []string
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.Mutex
}

func NewMetricsManager(cfg *config.AppConfig, logger logger.Logger) (*MetricsManager, error) {
	client, err := statsd.New(
		fmt.Sprintf("%s:%s", cfg.DDAgentHost, "8125"), // todo ポート確認
		statsd.WithNamespace("app."),                  // todo
		statsd.WithTags([]string{
			fmt.Sprintf("env:%s", cfg.Env),
			fmt.Sprintf("version:%s", cfg.Version),
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create StatsD client: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MetricsManager{
		logger: logger,
		client: client,
		prefix: "app.",
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (m *MetricsManager) Start() {
	go m.collectSystemMetrics()
}

func (m *MetricsManager) Stop() {
	m.cancel()
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.client != nil {
		if err := m.client.Close(); err != nil {
			m.logger.Error("Error closing StatsD client", "error", err)
		}
	}
}

func (m *MetricsManager) collectSystemMetrics() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			if m.client == nil {
				m.mu.Unlock()
				return
			}

			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)

			_ = m.client.Gauge("system.goroutines", float64(runtime.NumGoroutine()), nil, 1)
			_ = m.client.Gauge("system.memory.allocated", float64(mem.Alloc), []string{"type:heap"}, 1)
			_ = m.client.Gauge("system.memory.total", float64(mem.TotalAlloc), []string{"type:total"}, 1)

			m.mu.Unlock()

		case <-m.ctx.Done():
			return
		}
	}
}

func (m *MetricsManager) RecordHTTPMetrics(method, path string, statusCode int, duration time.Duration, responseSize int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client == nil {
		return
	}

	tags := []string{
		fmt.Sprintf("method:%s", method),
		fmt.Sprintf("path:%s", path),
		fmt.Sprintf("status:%d", statusCode),
	}

	_ = m.client.Incr("http.request.count", tags, 1)
	_ = m.client.Timing("http.request.duration", duration, tags, 1)
	_ = m.client.Histogram("http.response.size", float64(responseSize), tags, 1)

	if statusCode >= 400 {
		_ = m.client.Incr("http.request.errors", tags, 1)
	}
}
