BEGIN;

CREATE TABLE IF NOT EXISTS dish_views (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    viewed_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS dish_of_day (
    id boolean PRIMARY KEY DEFAULT true CHECK (id),
    dish_id bigint NOT NULL REFERENCES dishes(id) ON DELETE CASCADE,
    selected_at timestamptz NOT NULL DEFAULT now()
);

COMMIT;
