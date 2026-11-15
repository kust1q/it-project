package showcase

import (
	"context"
	"encoding/json"

	"project/internal/domain/entity"
)

type Repository interface {
	Restaurant(context.Context) (json.RawMessage, error)
	DishOfDay(context.Context) (*entity.Dish, error)
	Featured(context.Context, int64) ([]entity.Dish, error)
	New(context.Context, int64) ([]entity.Dish, error)
}
