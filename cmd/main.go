package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

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

	logger.Info("application started")
	<-ctx.Done()
	logger.Info("application stopping")
	return nil
}
