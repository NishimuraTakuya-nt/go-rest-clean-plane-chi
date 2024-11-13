package datadog

import (
	"fmt"

	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type Tracer struct {
	logger logger.Logger
	cfg    *config.AppConfig
}

func NewTracer(cfg *config.AppConfig, logger logger.Logger) *Tracer {
	return &Tracer{
		logger: logger,
		cfg:    cfg,
	}
}

func (t *Tracer) Start() error {
	if !t.cfg.DDEnabled {
		t.logger.Info("Datadog tracing is disabled")
		return nil
	}

	t.logger.Info("Initializing Datadog tracer")

	tracer.Start(
		tracer.WithService("go-rest-clean-plane-chi"),
		tracer.WithEnv(t.cfg.Env),
		tracer.WithServiceVersion(t.cfg.Version),
		tracer.WithAgentAddr(fmt.Sprintf("%s:%s",
			t.cfg.DDAgentHost,
			t.cfg.DDAgentPort,
		)),
		// サンプリングレートの設定
		tracer.WithSamplingRules([]tracer.SamplingRule{
			{
				Rate: 1.0,
			},
		}),
		// デバッグモードの設定
		tracer.WithDebugMode(t.cfg.LogLevel == "debug"),
		// プロファイラーの設定
		tracer.WithRuntimeMetrics(),
	)

	t.logger.Info("Datadog tracer initialized successfully",
		"agent_host", t.cfg.DDAgentHost,
		"agent_port", t.cfg.DDAgentPort,
		"env", t.cfg.Env,
	)
	return nil
}

func (t *Tracer) Stop() {
	if t.cfg.DDEnabled {
		t.logger.Info("Stopping Datadog tracer")
		tracer.Stop()
	}
}
