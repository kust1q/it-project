BEGIN;

DROP VIEW IF EXISTS dish_nutrition_view;

ALTER TABLE dish_nutrition ADD COLUMN IF NOT EXISTS calories      numeric(10,2);
ALTER TABLE dish_nutrition ADD COLUMN IF NOT EXISTS protein       numeric(10,2);
ALTER TABLE dish_nutrition ADD COLUMN IF NOT EXISTS fat           numeric(10,2);
ALTER TABLE dish_nutrition ADD COLUMN IF NOT EXISTS carbohydrates numeric(10,2);

UPDATE dish_nutrition dn
SET calories      = m.calories,
    protein       = m.protein,
    fat           = m.fat,
    carbohydrates = m.carbohydrates
FROM dish_nutrition_manual m
WHERE m.dish_id = dn.dish_id AND dn.is_manual;

DROP TABLE IF EXISTS dish_nutrition_manual;

ALTER TABLE dish_nutrition
    ALTER COLUMN is_manual SET DEFAULT false;

ALTER TABLE dish_nutrition ADD CONSTRAINT dish_nutrition_check CHECK (
    (NOT is_manual AND calories IS NULL AND protein IS NULL AND fat IS NULL AND carbohydrates IS NULL)
    OR
    (is_manual AND calories IS NOT NULL AND protein IS NOT NULL AND fat IS NOT NULL AND carbohydrates IS NOT NULL
     AND calories >= 0 AND protein >= 0 AND fat >= 0 AND carbohydrates >= 0)
);

COMMIT;