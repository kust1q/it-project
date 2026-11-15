package menu

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"project/internal/domain/entity"
	"project/internal/domain/errs"
)

func (s *Service) Categories(ctx context.Context) ([]entity.Category, error) {
	categories, err := s.repo.Categories(ctx)
	if err != nil {
		return nil, fmt.Errorf("menu.Categories.get: %w", err)
	}
	return categories, nil
}

func (s *Service) Category(ctx context.Context, id int64) (entity.Category, error) {
	category, err := s.repo.Category(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return category, fmt.Errorf("menu.Category.get: %w", errs.ErrNotFound)
	}
	if err != nil {
		return category, fmt.Errorf("menu.Category.get: %w", err)
	}
	return category, nil
}

func (s *Service) Dishes(ctx context.Context, filter entity.DishFilter) (entity.Page, error) {
	dishes, err := s.repo.Dishes(ctx, filter)
	if err != nil {
		return entity.Page{}, fmt.Errorf("menu.Dishes.get: %w", err)
	}
	return entity.Page{Items: dishes, Page: filter.Page, Limit: filter.Limit}, nil
}

func (s *Service) Dish(ctx context.Context, id int64) (entity.Dish, error) {
	dish, err := s.repo.Dish(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return dish, fmt.Errorf("menu.Dish.get: %w", errs.ErrNotFound)
	}
	if err != nil {
		return dish, fmt.Errorf("menu.Dish.get: %w", err)
	}
	if err := s.repo.RecordView(ctx, id); err != nil {
		s.logger.ErrorContext(ctx, "failed to record dish view", "dish_id", id, "error", err)
	}
	return dish, nil
}
