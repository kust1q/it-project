BEGIN;

DROP INDEX IF EXISTS one_active_menu_draft;

ALTER TABLE menu_draft_items ADD COLUMN IF NOT EXISTS original_version text;
ALTER TABLE dishes ADD COLUMN IF NOT EXISTS menu_version text;
ALTER TABLE restaurant ADD COLUMN IF NOT EXISTS opening_hours text;

DROP TABLE IF EXISTS restaurant_opening_hours;

COMMIT;