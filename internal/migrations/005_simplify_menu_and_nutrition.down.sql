BEGIN;

DROP VIEW IF EXISTS dish_nutrition_view;
CREATE TABLE IF NOT EXISTS dish_nutrition_manual (
    dish_id bigint PRIMARY KEY REFERENCES dishes(id) ON DELETE CASCADE,
    calories numeric(10,2) NOT NULL CHECK (calories >= 0),
    protein numeric(10,2) NOT NULL CHECK (protein >= 0),
    fat numeric(10,2) NOT NULL CHECK (fat >= 0),
    carbohydrates numeric(10,2) NOT NULL CHECK (carbohydrates >= 0)
);
ALTER TABLE dish_nutrition DROP CONSTRAINT IF EXISTS dish_nutrition_values_check;
ALTER TABLE dish_nutrition DROP COLUMN IF EXISTS calories, DROP COLUMN IF EXISTS protein,
    DROP COLUMN IF EXISTS fat, DROP COLUMN IF EXISTS carbohydrates;

CREATE OR REPLACE VIEW dish_nutrition_view AS
SELECT d.id AS dish_id, dn.is_manual,
    CASE WHEN dn.is_manual AND dnm.dish_id IS NOT NULL THEN dnm.calories ELSE t.calories END AS calories,
    CASE WHEN dn.is_manual AND dnm.dish_id IS NOT NULL THEN dnm.protein ELSE t.protein END AS protein,
    CASE WHEN dn.is_manual AND dnm.dish_id IS NOT NULL THEN dnm.fat ELSE t.fat END AS fat,
    CASE WHEN dn.is_manual AND dnm.dish_id IS NOT NULL THEN dnm.carbohydrates ELSE t.carbohydrates END AS carbohydrates
FROM dishes d
LEFT JOIN dish_nutrition dn ON dn.dish_id = d.id
LEFT JOIN dish_nutrition_manual dnm ON dnm.dish_id = d.id
LEFT JOIN LATERAL (
    SELECT SUM(di.weight * i.calories / 100) AS calories,
           SUM(di.weight * i.protein / 100) AS protein,
           SUM(di.weight * i.fat / 100) AS fat,
           SUM(di.weight * i.carbohydrates / 100) AS carbohydrates
    FROM dish_ingredients di JOIN ingredients i ON i.id = di.ingredient_id
    WHERE di.dish_id = d.id
) t ON true;

ALTER TABLE menu_drafts ADD COLUMN IF NOT EXISTS status text NOT NULL DEFAULT 'draft';
ALTER TABLE menu_drafts ADD CONSTRAINT menu_drafts_status_check CHECK (status IN ('draft', 'published'));
CREATE TABLE IF NOT EXISTS menu_published (
    draft_id bigint PRIMARY KEY REFERENCES menu_drafts(id) ON DELETE CASCADE,
    published_at timestamptz NOT NULL
);
ALTER TABLE menu_drafts DROP COLUMN IF EXISTS published_at;
COMMIT;
