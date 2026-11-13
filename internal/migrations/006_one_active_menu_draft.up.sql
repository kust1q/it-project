BEGIN;

CREATE UNIQUE INDEX IF NOT EXISTS one_active_menu_draft
ON menu_drafts ((true))
WHERE published_at IS NULL;

COMMIT;
