package handler

import (
	"context"
	"encoding/json"

	"project/internal/domain/entity"
)

type MenuService interface {
	Categories(context.Context) ([]entity.Category, error)
	Category(context.Context, int64) (entity.Category, error)
	Dishes(context.Context, entity.DishFilter) (entity.Page, error)
	Dish(context.Context, int64) (entity.Dish, error)
}

type ShowcaseService interface {
	Restaurant(context.Context) (json.RawMessage, error)
	Home(context.Context) (entity.Home, error)
}
