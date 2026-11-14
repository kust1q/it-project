BEGIN;

DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS promotion_dishes;
DROP TABLE IF EXISTS promotions;
DROP TABLE IF EXISTS admins;
DROP TABLE IF EXISTS dish_nutrition;
DROP TABLE IF EXISTS ingredient_allergens;
DROP TABLE IF EXISTS allergens;
DROP TABLE IF EXISTS dish_ingredients;
DROP TABLE IF EXISTS ingredients;
DROP TABLE IF EXISTS dish_photos;
DROP TABLE IF EXISTS dishes;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS restaurant_photos;
DROP TABLE IF EXISTS restaurant;

COMMIT;
