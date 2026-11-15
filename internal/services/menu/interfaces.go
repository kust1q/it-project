package menu

import (
	"context"

	"project/internal/domain/entity"
)

type Repository interface {
	Categories(context.Context) ([]entity.Category, error)
	Category(context.Context, int64) (entity.Category, error)
	Dishes(context.Context, entity.DishFilter) ([]entity.Dish, error)
	Dish(context.Context, int64) (entity.Dish, error)
	RecordView(context.Context, int64) error
}
