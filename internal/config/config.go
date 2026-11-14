package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.yaml.in/yaml/v3"
)

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return Config{}, fmt.Errorf("config.Load.loadEnv: %w", err)
	}

	file, err := os.Open("config/config.yml")
	if err != nil {
		return Config{}, fmt.Errorf("config.Load.openConfig: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, fmt.Errorf("config.Load.parseConfig: %w", err)
	}

	for _, item := range []struct {
		name string
		host *string
	}{
		{"POSTGRES_HOST", &cfg.Postgres.Host},
		{"REDIS_HOST", &cfg.Redis.Host},
		{"MINIO_HOST", &cfg.Minio.Host},
	} {
		if value := os.Getenv(item.name); value != "" {
			*item.host = value
		}
	}

	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")
	cfg.Redis.Password = os.Getenv("REDIS_PASSWORD")
	cfg.Minio.User = os.Getenv("MINIO_USER")
	cfg.Minio.Password = os.Getenv("MINIO_PASSWORD")

	var missing []string
	for _, item := range []struct{ name, value string }{
		{"app.port", cfg.App.Port}, {"postgres.host", cfg.Postgres.Host}, {"postgres.port", cfg.Postgres.Port},
		{"postgres.sslmode", cfg.Postgres.SSLMode}, {"redis.host", cfg.Redis.Host}, {"redis.port", cfg.Redis.Port},
		{"minio.host", cfg.Minio.Host}, {"minio.port", cfg.Minio.Port}, {"minio.bucket", cfg.Minio.BucketName},
		{"POSTGRES_USER", cfg.Postgres.User}, {"POSTGRES_PASSWORD", cfg.Postgres.Password}, {"POSTGRES_DB", cfg.Postgres.DBName},
		{"REDIS_PASSWORD", cfg.Redis.Password},
		{"MINIO_USER", cfg.Minio.User}, {"MINIO_PASSWORD", cfg.Minio.Password},
	} {
		if strings.TrimSpace(item.value) == "" {
			missing = append(missing, item.name)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("config.Load.validate: missing configuration: %s", strings.Join(missing, ", "))
	}
	if cfg.Admin.SessionTTL <= 0 {
		return Config{}, fmt.Errorf("config.Load.validate: admin.session_ttl must be a positive duration")
	}
	return cfg, nil
}
