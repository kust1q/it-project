BEGIN;

ALTER TABLE restaurant ADD COLUMN IF NOT EXISTS opening_hours text NOT NULL DEFAULT '';

DROP TABLE IF EXISTS restaurant_working_time;

COMMIT;