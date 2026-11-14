BEGIN;

CREATE TABLE IF NOT EXISTS dish_views (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE RESTRICT,
    viewed_at timestamptz NOT NULL DEFAULT now()
);

COMMIT;
