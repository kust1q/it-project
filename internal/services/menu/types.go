package menu

import "log/slog"

type Service struct {
	repo   Repository
	logger *slog.Logger
}
