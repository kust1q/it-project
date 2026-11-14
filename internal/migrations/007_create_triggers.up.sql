BEGIN;

CREATE OR REPLACE FUNCTION fn_create_dish_nutrition()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    INSERT INTO dish_nutrition (dish_id, is_manual)
    VALUES (NEW.id, false)
    ON CONFLICT (dish_id) DO NOTHING;

    RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_dishes_create_nutrition ON dishes;
CREATE TRIGGER trg_dishes_create_nutrition
AFTER INSERT ON dishes
FOR EACH ROW
EXECUTE FUNCTION fn_create_dish_nutrition();

CREATE OR REPLACE FUNCTION fn_audit_log()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
DECLARE
    v_admin_id  bigint;
    v_id_col    text := COALESCE(TG_ARGV[0], 'id');
    v_old_data  jsonb;
    v_new_data  jsonb;
    v_entity_id bigint;
BEGIN
    BEGIN
        v_admin_id := NULLIF(btrim(current_setting('app.current_admin_id', true)), '')::bigint;
    EXCEPTION WHEN OTHERS THEN
        v_admin_id := NULL;
    END;

    IF TG_OP = 'INSERT' THEN
        v_new_data  := to_jsonb(NEW);
        v_entity_id := (v_new_data ->> v_id_col)::bigint;
    ELSIF TG_OP = 'UPDATE' THEN
        IF NEW IS NOT DISTINCT FROM OLD THEN
            RETURN NULL;
        END IF;
        v_old_data  := to_jsonb(OLD);
        v_new_data  := to_jsonb(NEW);
        v_entity_id := (v_new_data ->> v_id_col)::bigint;
    ELSIF TG_OP = 'DELETE' THEN
        v_old_data  := to_jsonb(OLD);
        v_entity_id := (v_old_data ->> v_id_col)::bigint;
    END IF;

    INSERT INTO audit_log (admin_id, entity_type, entity_id, operation, old_data, new_data)
    VALUES (v_admin_id, TG_TABLE_NAME, v_entity_id, TG_OP, v_old_data, v_new_data);

    RETURN NULL;
END;
$$;

DROP TRIGGER IF EXISTS trg_dishes_audit ON dishes;
CREATE TRIGGER trg_dishes_audit
AFTER INSERT OR UPDATE OR DELETE ON dishes
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('id');

DROP TRIGGER IF EXISTS trg_ingredients_audit ON ingredients;
CREATE TRIGGER trg_ingredients_audit
AFTER INSERT OR UPDATE OR DELETE ON ingredients
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('id');

DROP TRIGGER IF EXISTS trg_dish_ingredients_audit ON dish_ingredients;
CREATE TRIGGER trg_dish_ingredients_audit
AFTER INSERT OR UPDATE OR DELETE ON dish_ingredients
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('dish_id');

DROP TRIGGER IF EXISTS trg_allergens_audit ON allergens;
CREATE TRIGGER trg_allergens_audit
AFTER INSERT OR UPDATE OR DELETE ON allergens
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('id');

DROP TRIGGER IF EXISTS trg_ingredient_allergens_audit ON ingredient_allergens;
CREATE TRIGGER trg_ingredient_allergens_audit
AFTER INSERT OR UPDATE OR DELETE ON ingredient_allergens
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('ingredient_id');

DROP TRIGGER IF EXISTS trg_promotions_audit ON promotions;
CREATE TRIGGER trg_promotions_audit
AFTER INSERT OR UPDATE OR DELETE ON promotions
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('id');

DROP TRIGGER IF EXISTS trg_promotion_dishes_audit ON promotion_dishes;
CREATE TRIGGER trg_promotion_dishes_audit
AFTER INSERT OR UPDATE OR DELETE ON promotion_dishes
FOR EACH ROW
EXECUTE FUNCTION fn_audit_log('promotion_id');

COMMIT;