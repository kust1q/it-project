BEGIN;

CREATE TABLE IF NOT EXISTS restaurant (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL CHECK (btrim(name) <> ''),
    description text NOT NULL DEFAULT '',
    address text NOT NULL CHECK (btrim(address) <> ''),
    opening_hours text NOT NULL DEFAULT '',
    phone text,
    email text
);

CREATE TABLE IF NOT EXISTS restaurant_photos (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id bigint NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
    object_key text NOT NULL UNIQUE,
    caption text NOT NULL DEFAULT '',
    position integer NOT NULL DEFAULT 0 CHECK (position >= 0)
);

CREATE TABLE IF NOT EXISTS categories (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    position integer NOT NULL DEFAULT 0 CHECK (position >= 0)
);

CREATE TABLE IF NOT EXISTS dishes (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    category_id bigint NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    name text NOT NULL CHECK (btrim(name) <> ''),
    description text NOT NULL DEFAULT '',
    portion_weight numeric(10,2) NOT NULL CHECK (portion_weight > 0),
    price numeric(12,2) NOT NULL CHECK (price > 0),
    is_available boolean NOT NULL DEFAULT false,
    is_stopped boolean NOT NULL DEFAULT false,
    is_archived boolean NOT NULL DEFAULT false,
    signature_since timestamptz
);

CREATE TABLE IF NOT EXISTS dish_photos (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    object_key text NOT NULL UNIQUE,
    position integer NOT NULL DEFAULT 0 CHECK (position >= 0),
    is_primary boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX IF NOT EXISTS dish_photos_one_primary ON dish_photos (dish_id) WHERE is_primary;

CREATE TABLE IF NOT EXISTS ingredients (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE CHECK (btrim(name) <> ''),
    calories numeric(10,2) NOT NULL CHECK (calories >= 0),
    protein numeric(10,2) NOT NULL CHECK (protein >= 0),
    fat numeric(10,2) NOT NULL CHECK (fat >= 0),
    carbohydrates numeric(10,2) NOT NULL CHECK (carbohydrates >= 0)
);

CREATE TABLE IF NOT EXISTS dish_ingredients (
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    ingredient_id bigint NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    weight numeric(10,2) NOT NULL CHECK (weight > 0),
    PRIMARY KEY (dish_id, ingredient_id)
);

CREATE TABLE IF NOT EXISTS allergens (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL UNIQUE CHECK (btrim(name) <> '')
);

CREATE TABLE IF NOT EXISTS ingredient_allergens (
    ingredient_id bigint NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    allergen_id bigint NOT NULL REFERENCES allergens(id) ON DELETE RESTRICT,
    PRIMARY KEY (ingredient_id, allergen_id)
);

CREATE TABLE IF NOT EXISTS dish_nutrition (
    dish_id bigint PRIMARY KEY REFERENCES dishes(id) ON DELETE CASCADE,
    is_manual boolean NOT NULL DEFAULT false,
    calories numeric(10,2),
    protein numeric(10,2),
    fat numeric(10,2),
    carbohydrates numeric(10,2),
    CHECK (
        (NOT is_manual AND calories IS NULL AND protein IS NULL AND fat IS NULL AND carbohydrates IS NULL)
        OR
        (is_manual AND calories IS NOT NULL AND protein IS NOT NULL AND fat IS NOT NULL AND carbohydrates IS NOT NULL
         AND calories >= 0 AND protein >= 0 AND fat >= 0 AND carbohydrates >= 0)
    )
);

CREATE TABLE IF NOT EXISTS admins (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    login text NOT NULL UNIQUE CHECK (btrim(login) <> ''),
    password_hash text NOT NULL CHECK (btrim(password_hash) <> ''),
    is_active boolean NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS promotions (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name text NOT NULL CHECK (btrim(name) <> ''),
    discount_percent numeric(5,2) NOT NULL CHECK (discount_percent > 0 AND discount_percent < 100),
    starts_at timestamptz NOT NULL,
    ends_at timestamptz NOT NULL,
    is_enabled boolean NOT NULL DEFAULT false,
    CHECK (starts_at < ends_at)
);

CREATE TABLE IF NOT EXISTS promotion_dishes (
    promotion_id bigint NOT NULL REFERENCES promotions(id) ON DELETE CASCADE,
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE RESTRICT,
    PRIMARY KEY (promotion_id, dish_id)
);

CREATE TABLE IF NOT EXISTS audit_log (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    admin_id bigint REFERENCES admins(id) ON DELETE RESTRICT,
    entity_type text NOT NULL,
    entity_id bigint NOT NULL,
    operation text NOT NULL CHECK (operation IN ('INSERT', 'UPDATE', 'DELETE')),
    changed_at timestamptz NOT NULL DEFAULT now(),
    old_data jsonb,
    new_data jsonb
);

COMMIT;
