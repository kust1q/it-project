package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"project/internal/domain/entity"
)

type MenuRepository struct{ db *sqlx.DB }

func NewMenuRepository(db *sqlx.DB) *MenuRepository { return &MenuRepository{db: db} }

// Prices are multiplied as numeric values in PostgreSQL and rounded only once.
const dishJSON = `jsonb_build_object(
'id', d.id, 'category_id', d.category_id, 'name', d.name, 'description', d.description,
'portion_weight', d.portion_weight, 'price', d.price, 'final_price', COALESCE(discounts.final_price, d.price),
'is_signature', d.signature_since IS NOT NULL,
'photos', COALESCE((SELECT jsonb_agg(jsonb_build_object('id', ph.id, 'object_key', ph.object_key, 'position', ph.position, 'is_primary', ph.is_primary) ORDER BY ph.position, ph.id) FROM dish_photos ph WHERE ph.dish_id=d.id), '[]'::jsonb),
'ingredients', COALESCE((SELECT jsonb_agg(jsonb_build_object('id', i.id, 'name', i.name, 'weight', di.weight) ORDER BY i.name, i.id) FROM dish_ingredients di JOIN ingredients i ON i.id=di.ingredient_id WHERE di.dish_id=d.id), '[]'::jsonb),
'allergens', COALESCE((SELECT jsonb_agg(jsonb_build_object('id', a.id, 'name', a.name) ORDER BY a.name, a.id) FROM allergens a WHERE EXISTS (SELECT 1 FROM dish_ingredients di JOIN ingredient_allergens ia ON ia.ingredient_id=di.ingredient_id WHERE di.dish_id=d.id AND ia.allergen_id=a.id)), '[]'::jsonb),
'nutrition', (SELECT jsonb_build_object('is_manual', n.is_manual, 'calories', n.calories, 'protein', n.protein, 'fat', n.fat, 'carbohydrates', n.carbohydrates) FROM dish_nutrition_view n WHERE n.dish_id=d.id),
'promotions', COALESCE(discounts.promotions, '[]'::jsonb))`

const discountJoin = ` LEFT JOIN LATERAL (
WITH RECURSIVE active AS (
 SELECT p.id, p.name, p.discount_percent, row_number() OVER (ORDER BY p.id) AS rn
 FROM promotions p WHERE p.is_enabled AND p.starts_at <= now() AND p.ends_at > now()
 AND (p.applies_to_all OR EXISTS (SELECT 1 FROM promotion_dishes pd WHERE pd.promotion_id=p.id AND pd.dish_id=d.id))
), calc AS (
 SELECT 0::bigint AS rn, d.price::numeric AS value
 UNION ALL SELECT a.rn, c.value * (100-a.discount_percent)/100 FROM calc c JOIN active a ON a.rn=c.rn+1
)
SELECT (SELECT round(value, 2) FROM calc ORDER BY rn DESC LIMIT 1) AS final_price,
 (SELECT jsonb_agg(jsonb_build_object('id',id,'name',name,'discount_percent',discount_percent) ORDER BY id) FROM active) AS promotions
) discounts ON true `
const visible = `NOT d.is_archived AND (d.stopped_until IS NULL OR d.stopped_until <= now())`

func decodeDish(raw []byte) (entity.Dish, error) {
	var d entity.Dish
	err := json.Unmarshal(raw, &d)
	return d, err
}
func (r *MenuRepository) Categories(ctx context.Context) ([]entity.Category, error) {
	var result []entity.Category
	err := r.db.SelectContext(ctx, &result, `SELECT id,name,position FROM categories ORDER BY position,id`)
	if err != nil {
		return nil, fmt.Errorf("postgres.Categories.select: %w", err)
	}
	if result == nil {
		result = []entity.Category{}
	}
	return result, nil
}
func (r *MenuRepository) Category(ctx context.Context, id int64) (entity.Category, error) {
	var c entity.Category
	err := r.db.GetContext(ctx, &c, `SELECT id,name,position FROM categories WHERE id=$1`, id)
	if err != nil {
		return c, fmt.Errorf("postgres.Category.select: %w", err)
	}
	return c, nil
}
func (r *MenuRepository) list(ctx context.Context, where, order string, args ...any) ([]entity.Dish, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT (`+dishJSON+`) - 'ingredients' - 'allergens' - 'nutrition' FROM dishes d `+discountJoin+` WHERE `+visible+` AND `+where+` ORDER BY `+order, args...)
	if err != nil {
		return nil, fmt.Errorf("postgres.list.query: %w", err)
	}
	defer rows.Close()
	result := []entity.Dish{}
	for rows.Next() {
		var raw []byte
		if err = rows.Scan(&raw); err != nil {
			return nil, fmt.Errorf("postgres.list.scan: %w", err)
		}
		d, e := decodeDish(raw)
		if e != nil {
			return nil, fmt.Errorf("postgres.list.decode: %w", e)
		}
		result = append(result, d)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("postgres.list.rows: %w", err)
	}
	return result, nil
}
func (r *MenuRepository) Dishes(ctx context.Context, f entity.DishFilter) ([]entity.Dish, error) {
	order := map[string]string{"name": "d.name ASC,d.id ASC", "-name": "d.name DESC,d.id ASC", "price": "discounts.final_price ASC,d.id ASC", "-price": "discounts.final_price DESC,d.id ASC", "newest": "d.created_at DESC,d.id ASC"}[f.Sort]
	if order == "" {
		order = "d.id ASC"
	}
	return r.list(ctx, `($1::bigint=0 OR d.category_id=$1) AND ($2='' OR d.name ILIKE '%' || $2 || '%')`, order+` LIMIT $3 OFFSET $4`, f.CategoryID, strings.TrimSpace(f.Search), f.Limit, (f.Page-1)*f.Limit)
}
func (r *MenuRepository) Dish(ctx context.Context, id int64) (entity.Dish, error) {
	var raw []byte
	err := r.db.QueryRowContext(ctx, `SELECT `+dishJSON+` FROM dishes d `+discountJoin+` WHERE d.id=$1 AND `+visible, id).Scan(&raw)
	if err != nil {
		return entity.Dish{}, fmt.Errorf("postgres.Dish.select: %w", err)
	}
	d, err := decodeDish(raw)
	if err != nil {
		return d, fmt.Errorf("postgres.Dish.decode: %w", err)
	}
	return d, nil
}
func (r *MenuRepository) RecordView(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO dish_views(dish_id) SELECT id FROM dishes d WHERE id=$1 AND `+visible, id)
	if err != nil {
		return fmt.Errorf("postgres.RecordView.insert: %w", err)
	}
	return nil
}
func (r *MenuRepository) Restaurant(ctx context.Context) (json.RawMessage, error) {
	var raw []byte
	err := r.db.QueryRowContext(ctx, `SELECT jsonb_build_object('id',r.id,'name',r.name,'description',r.description,'address',r.address,
 'contacts',COALESCE((SELECT jsonb_agg(jsonb_build_object('id',c.id,'type',c.type,'value',c.value) ORDER BY c.id) FROM restaurant_contacts c WHERE c.restaurant_id=r.id),'[]'::jsonb),
 'working_time',COALESCE((SELECT jsonb_agg(jsonb_build_object('day_of_week',w.day_of_week,'open_time',w.open_time,'close_time',w.close_time) ORDER BY w.day_of_week) FROM restaurant_working_time w WHERE w.restaurant_id=r.id),'[]'::jsonb),
 'photos',COALESCE((SELECT jsonb_agg(jsonb_build_object('id',p.id,'object_key',p.object_key,'caption',p.caption,'position',p.position) ORDER BY p.position,p.id) FROM restaurant_photos p WHERE p.restaurant_id=r.id),'[]'::jsonb)) FROM restaurant r ORDER BY r.id LIMIT 1`).Scan(&raw)
	if err != nil {
		return nil, fmt.Errorf("postgres.Restaurant.select: %w", err)
	}
	return json.RawMessage(raw), nil
}

func (r *MenuRepository) DishOfDay(ctx context.Context) (*entity.Dish, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("postgres.DishOfDay.begin: %w", err)
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(73194821)`); err != nil {
		return nil, fmt.Errorf("postgres.DishOfDay.lock: %w", err)
	}
	// A single database timestamp determines both the Moscow day and the one-hour window.
	var current sql.NullInt64
	err = tx.QueryRowContext(ctx, `WITH bounds AS (SELECT now() AS now_at, date_trunc('day', now() AT TIME ZONE 'Europe/Moscow') AT TIME ZONE 'Europe/Moscow' AS day_start),
 first_view AS (SELECT min(viewed_at) AS at FROM dish_views,bounds WHERE viewed_at>=day_start AND viewed_at<=now_at)
 SELECT dd.dish_id FROM dish_of_day dd JOIN dishes d ON d.id=dd.dish_id CROSS JOIN bounds b CROSS JOIN first_view f
 WHERE dd.selected_at>=b.day_start AND dd.selected_at + interval '1 hour'>b.now_at
 AND f.at IS NOT NULL AND f.at + interval '1 hour'<=b.now_at AND `+visible).Scan(&current)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("postgres.DishOfDay.current: %w", err)
	}
	if !current.Valid {
		var candidate int64
		err = tx.QueryRowContext(ctx, `WITH bounds AS (SELECT now() AS now_at,date_trunc('day',now() AT TIME ZONE 'Europe/Moscow') AT TIME ZONE 'Europe/Moscow' AS day_start),
 first_view AS (SELECT min(viewed_at) AS at FROM dish_views,bounds WHERE viewed_at>=day_start AND viewed_at<=now_at)
 SELECT d.id FROM dishes d JOIN dish_views v ON v.dish_id=d.id CROSS JOIN bounds b CROSS JOIN first_view f
 WHERE `+visible+` AND v.viewed_at>=b.day_start AND v.viewed_at<=b.now_at AND f.at + interval '1 hour'<=b.now_at
 GROUP BY d.id ORDER BY count(*) DESC,max(v.viewed_at) DESC,max(v.id) DESC,d.id ASC LIMIT 1`).Scan(&candidate)
		if errors.Is(err, sql.ErrNoRows) {
			if _, err = tx.ExecContext(ctx, `DELETE FROM dish_of_day`); err != nil {
				return nil, fmt.Errorf("postgres.DishOfDay.clear: %w", err)
			}
			if err = tx.Commit(); err != nil {
				return nil, fmt.Errorf("postgres.DishOfDay.commit: %w", err)
			}
			return nil, nil
		}
		if err != nil {
			return nil, fmt.Errorf("postgres.DishOfDay.candidate: %w", err)
		}
		// Keeping selected_at when the leader stays the same avoids extending the lock-in period.
		_, err = tx.ExecContext(ctx, `INSERT INTO dish_of_day(dish_id,selected_at) VALUES($1,now()) ON CONFLICT(id) DO UPDATE SET dish_id=EXCLUDED.dish_id,selected_at=CASE WHEN dish_of_day.dish_id=EXCLUDED.dish_id AND dish_of_day.selected_at >= (date_trunc('day',now() AT TIME ZONE 'Europe/Moscow') AT TIME ZONE 'Europe/Moscow') THEN dish_of_day.selected_at ELSE EXCLUDED.selected_at END`, candidate)
		if err != nil {
			return nil, fmt.Errorf("postgres.DishOfDay.save: %w", err)
		}
		current = sql.NullInt64{Int64: candidate, Valid: true}
	}
	var raw []byte
	err = tx.QueryRowContext(ctx, `SELECT `+dishJSON+` FROM dishes d `+discountJoin+` WHERE d.id=$1 AND `+visible, current.Int64).Scan(&raw)
	if errors.Is(err, sql.ErrNoRows) {
		if err = tx.Commit(); err != nil {
			return nil, fmt.Errorf("postgres.DishOfDay.commit: %w", err)
		}
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("postgres.DishOfDay.select: %w", err)
	}
	d, err := decodeDish(raw)
	if err != nil {
		return nil, fmt.Errorf("postgres.DishOfDay.decode: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("postgres.DishOfDay.commit: %w", err)
	}
	return &d, nil
}
func (r *MenuRepository) Featured(ctx context.Context, exclude int64) ([]entity.Dish, error) {
	return r.list(ctx, `d.signature_since IS NOT NULL AND d.id<>$1`, `d.signature_since ASC,d.id ASC`, exclude)
}
func (r *MenuRepository) New(ctx context.Context, exclude int64) ([]entity.Dish, error) {
	return r.list(ctx, `d.signature_since IS NULL AND d.created_at > now()-interval '3 months' AND d.id<>$1`, `d.created_at DESC,d.id ASC`, exclude)
}
