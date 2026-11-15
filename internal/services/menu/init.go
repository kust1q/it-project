package menu

import "log/slog"

func New(repo Repository, logger *slog.Logger) *Service {
	return &Service{repo: repo, logger: logger}
}
