package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/docs/swagger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/config"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/logger"
	"github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/infrastructure/telemetry/datadog"
)

// @title Go REST Clean API with Chi
// @version 1.0
// @description This is a sample server for a Go REST API using clean architecture.
// @host localhost:8081
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	configLoader := config.NewLoader()
	cfg, err := configLoader.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	logger := logger.NewLogger(cfg)

	// metrics
	metricsManager, err := datadog.NewMetricsManager(cfg, logger)
	if err != nil {
		logger.Error("Failed to initialize metrics manager", "error", err)
		return fmt.Errorf("failed to initialize metrics manager: %w", err)
	}
	metricsManager.Start()
	defer metricsManager.Stop()
	// tracer
	ddTracer := datadog.NewTracer(cfg, logger)
	if err := ddTracer.Start(); err != nil {
		logger.Error("Failed to initialize Datadog tracer", "error", err)
		return err
	}
	defer ddTracer.Stop()

	router, err := InitializeRouter(cfg, logger, metricsManager)
	if err != nil {
		return err
	}
	h := router.Setup()

	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: h,
	}

	// シグナルを受け取るためのコンテキストを設定
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt, os.Kill)
	defer stop()

	// サーバーをゴルーチンで起動
	go func() {
		logger.Info("Server started", slog.String("address", cfg.ServerAddress))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Server listen failed", slog.String("error", err.Error()))
		}
	}()

	// シグナルを待機
	<-ctx.Done()
	logger.Info("Shutdown signal received")

	// シャットダウンのためのコンテキストを作成
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 順番にシャットダウン
	// 1. HTTPサーバー
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Server forced to shutdown", slog.String("error", err.Error()))
		return err
	}
	// 2. Telemetry
	//telemetryCleanup()

	// その他のリソースのクリーンアップ
	//if err := graphQLClient.Close(); err != nil {
	//	logger.Error("Error closing GraphQL client", slog.String("error", err.Error()))
	//}

	logger.Info("Server exited properly")
	return nil
}
