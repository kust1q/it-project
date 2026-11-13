BEGIN;

ALTER TABLE menu_drafts ADD COLUMN IF NOT EXISTS published_at timestamptz;
DROP TABLE IF EXISTS menu_published;
ALTER TABLE menu_drafts DROP COLUMN IF EXISTS status;

DROP VIEW IF EXISTS dish_nutrition_full;
ALTER TABLE dish_nutrition
    ADD COLUMN IF NOT EXISTS calories numeric(10,2),
    ADD COLUMN IF NOT EXISTS protein numeric(10,2),
    ADD COLUMN IF NOT EXISTS fat numeric(10,2),
    ADD COLUMN IF NOT EXISTS carbohydrates numeric(10,2);

ALTER TABLE dish_nutrition ADD CONSTRAINT dish_nutrition_values_check CHECK (
    (NOT is_manual AND calories IS NULL AND protein IS NULL AND fat IS NULL AND carbohydrates IS NULL)
    OR
    (calories IS NOT NULL AND protein IS NOT NULL AND fat IS NOT NULL AND carbohydrates IS NOT NULL
     AND calories >= 0 AND protein >= 0 AND fat >= 0 AND carbohydrates >= 0)
);
DROP TABLE IF EXISTS dish_nutrition_manual;

CREATE OR REPLACE VIEW dish_nutrition_full AS
SELECT d.id AS dish_id, dn.is_manual,
    CASE WHEN dn.is_manual THEN dn.calories ELSE t.calories END AS calories,
    CASE WHEN dn.is_manual THEN dn.protein ELSE t.protein END AS protein,
    CASE WHEN dn.is_manual THEN dn.fat ELSE t.fat END AS fat,
    CASE WHEN dn.is_manual THEN dn.carbohydrates ELSE t.carbohydrates END AS carbohydrates
FROM dishes d
LEFT JOIN dish_nutrition dn ON dn.dish_id = d.id
LEFT JOIN LATERAL (
    SELECT SUM(di.weight * i.calories / 100) AS calories,
           SUM(di.weight * i.protein / 100) AS protein,
           SUM(di.weight * i.fat / 100) AS fat,
           SUM(di.weight * i.carbohydrates / 100) AS carbohydrates
    FROM dish_ingredients di JOIN ingredients i ON i.id = di.ingredient_id
    WHERE di.dish_id = d.id
) t ON true;

COMMIT;
