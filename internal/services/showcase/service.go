package showcase

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"project/internal/domain/entity"
	"project/internal/domain/errs"
)

func (s *Service) Restaurant(ctx context.Context) (json.RawMessage, error) {
	restaurant, err := s.repo.Restaurant(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("showcase.Restaurant.get: %w", errs.ErrNotFound)
	}
	if err != nil {
		return nil, fmt.Errorf("showcase.Restaurant.get: %w", err)
	}
	return restaurant, nil
}

func (s *Service) Home(ctx context.Context) (entity.Home, error) {
	restaurant, err := s.Restaurant(ctx)
	if err != nil {
		return entity.Home{}, fmt.Errorf("showcase.Home.restaurant: %w", err)
	}
	day, err := s.repo.DishOfDay(ctx)
	if err != nil {
		return entity.Home{}, fmt.Errorf("showcase.Home.dishOfDay: %w", err)
	}
	var exclude int64
	if day != nil {
		exclude = day.ID
	}
	featured, err := s.repo.Featured(ctx, exclude)
	if err != nil {
		return entity.Home{}, fmt.Errorf("showcase.Home.featured: %w", err)
	}
	fresh, err := s.repo.New(ctx, exclude)
	if err != nil {
		return entity.Home{}, fmt.Errorf("showcase.Home.new: %w", err)
	}
	return entity.Home{Restaurant: restaurant, DishOfDay: day, SignatureDishes: featured, NewDishes: fresh}, nil
}
