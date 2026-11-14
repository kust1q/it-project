BEGIN;

CREATE TABLE IF NOT EXISTS restaurant_working_time (
    id bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    restaurant_id bigint NOT NULL REFERENCES restaurant(id) ON DELETE CASCADE,
    day_of_week smallint NOT NULL CHECK (day_of_week BETWEEN 1 AND 7),
    open_time time NOT NULL,
    close_time time NOT NULL,
    CHECK (open_time < close_time),
    UNIQUE (restaurant_id, day_of_week)
);

ALTER TABLE restaurant DROP COLUMN IF EXISTS opening_hours;

COMMIT;