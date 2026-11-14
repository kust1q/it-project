BEGIN;

DROP TRIGGER IF EXISTS trg_promotion_dishes_audit ON promotion_dishes;
DROP TRIGGER IF EXISTS trg_promotions_audit ON promotions;
DROP TRIGGER IF EXISTS trg_ingredient_allergens_audit ON ingredient_allergens;
DROP TRIGGER IF EXISTS trg_allergens_audit ON allergens;
DROP TRIGGER IF EXISTS trg_dish_ingredients_audit ON dish_ingredients;
DROP TRIGGER IF EXISTS trg_ingredients_audit ON ingredients;
DROP TRIGGER IF EXISTS trg_dishes_audit ON dishes;
DROP TRIGGER IF EXISTS trg_dishes_create_nutrition ON dishes;

DROP FUNCTION IF EXISTS fn_audit_log();
DROP FUNCTION IF EXISTS fn_create_dish_nutrition();

COMMIT;