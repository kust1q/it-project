package handler

import "log/slog"

type Handler struct {
	menu     MenuService
	showcase ShowcaseService
	logger   *slog.Logger
}
