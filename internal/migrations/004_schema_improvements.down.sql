BEGIN;

DROP INDEX IF EXISTS one_active_menu_draft;

ALTER TABLE menu_draft_items ADD COLUMN IF NOT EXISTS original_version bigint NOT NULL DEFAULT 1 CHECK (original_version > 0);
ALTER TABLE dishes ADD COLUMN IF NOT EXISTS menu_version bigint NOT NULL DEFAULT 1 CHECK (menu_version > 0);
ALTER TABLE restaurant ADD COLUMN IF NOT EXISTS opening_hours text NOT NULL DEFAULT '';

DROP TABLE IF EXISTS restaurant_opening_hours;

COMMIT;