-- recipes 表增加浏览量和收藏数字段
ALTER TABLE recipes ADD COLUMN IF NOT EXISTS view_count INT NOT NULL DEFAULT 0;
ALTER TABLE recipes ADD COLUMN IF NOT EXISTS favorite_count INT NOT NULL DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_recipes_view_count ON recipes(view_count DESC);
CREATE INDEX IF NOT EXISTS idx_recipes_favorite_count ON recipes(favorite_count DESC);

-- 用户自创菜谱表
CREATE TABLE IF NOT EXISTS user_recipes (
  id BIGSERIAL PRIMARY KEY,
  user_openid VARCHAR(64) NOT NULL,
  dish_name VARCHAR(255) NOT NULL,
  servings INT NOT NULL DEFAULT 2,
  recipe_json JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_user_recipes_user_openid ON user_recipes(user_openid);
CREATE INDEX IF NOT EXISTS idx_user_recipes_created_at ON user_recipes(created_at DESC);

-- 收藏表
CREATE TABLE IF NOT EXISTS favorites (
  id BIGSERIAL PRIMARY KEY,
  user_openid VARCHAR(64) NOT NULL,
  recipe_id BIGINT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(user_openid, recipe_id)
);
CREATE INDEX IF NOT EXISTS idx_favorites_user_openid ON favorites(user_openid);
CREATE INDEX IF NOT EXISTS idx_favorites_recipe_id ON favorites(recipe_id);
