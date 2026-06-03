DROP TABLE IF EXISTS favorites;
DROP TABLE IF EXISTS user_recipes;

DROP INDEX IF EXISTS idx_recipes_view_count;
DROP INDEX IF EXISTS idx_recipes_favorite_count;
ALTER TABLE recipes DROP COLUMN IF EXISTS view_count;
ALTER TABLE recipes DROP COLUMN IF EXISTS favorite_count;
