package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"project/internal/api/handler"
	"project/internal/config"
	"project/internal/repository/clickhouse"
	"project/internal/repository/minio"
	"project/internal/repository/postgres"
	"project/internal/repository/rabbitmq"
	"project/internal/repository/redis"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := run(logger); err != nil {
		logger.Error("application failed", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("main.run.loadConfig: %w", err)
	}

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		return fmt.Errorf("main.run.connectPostgres: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Error("failed to close postgres", "error", err)
		}
	}()
	logger.Info("connected to postgres")

	redisClient, err := redis.New(&cfg.Redis)
	if err != nil {
		return fmt.Errorf("main.run.connectRedis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			logger.Error("failed to close redis", "error", err)
		}
	}()
	logger.Info("connected to redis")

	rabbitConn, err := rabbitmq.New(&cfg.RabbitMQ)
	if err != nil {
		return fmt.Errorf("main.run.connectRabbitMQ: %w", err)
	}
	defer func() {
		if err := rabbitConn.Close(); err != nil {
			logger.Error("failed to close rabbitmq", "error", err)
		}
	}()
	logger.Info("connected to rabbitmq")

	clickhouseConn, err := clickhouse.New(&cfg.ClickHouse)
	if err != nil {
		return fmt.Errorf("main.run.connectClickHouse: %w", err)
	}
	defer func() {
		if err := clickhouseConn.Close(); err != nil {
			logger.Error("failed to close clickhouse", "error", err)
		}
	}()
	logger.Info("connected to clickhouse")

	_, err = minio.New(&cfg.Minio)
	if err != nil {
		return fmt.Errorf("main.run.connectMinio: %w", err)
	}
	logger.Info("connected to minio")

	server := &http.Server{
		Addr:              net.JoinHostPort("", cfg.App.Port),
		Handler:           handler.New().InitRouters(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		return fmt.Errorf("main.run.listenHTTP: %w", err)
	}

	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.Serve(listener)
	}()
	logger.Info("http server started", "address", server.Addr)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("main.run.serveHTTP: %w", err)
		}
		return nil
	case <-ctx.Done():
		logger.Info("application stopping")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		if closeErr := server.Close(); closeErr != nil {
			logger.Error("failed to close http server", "error", closeErr)
		}
		return fmt.Errorf("main.run.shutdownHTTP: %w", err)
	}
	if err := <-serverErrors; !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("main.run.serveHTTP: %w", err)
	}
	logger.Info("http server stopped")
	return nil
}
