CREATE TABLE IF NOT EXISTS analyze_tasks (
  id BIGSERIAL PRIMARY KEY,
  task_id VARCHAR(64) UNIQUE NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'pending',
  message TEXT,
  recipe_id BIGINT,
  error_message TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_analyze_tasks_task_id ON analyze_tasks(task_id);
CREATE INDEX IF NOT EXISTS idx_analyze_tasks_status ON analyze_tasks(status);

CREATE TABLE IF NOT EXISTS bilibili_videos (
  id BIGSERIAL PRIMARY KEY,
  bvid VARCHAR(64) UNIQUE,
  aid BIGINT,
  cid BIGINT,
  title TEXT,
  description TEXT,
  owner_name VARCHAR(255),
  duration_seconds INT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_bilibili_videos_bvid ON bilibili_videos(bvid);
CREATE INDEX IF NOT EXISTS idx_bilibili_videos_aid ON bilibili_videos(aid);

CREATE TABLE IF NOT EXISTS video_text_sources (
  id BIGSERIAL PRIMARY KEY,
  video_id BIGINT REFERENCES bilibili_videos(id) ON DELETE CASCADE,
  source_type VARCHAR(32) NOT NULL,
  content TEXT,
  raw_json JSONB,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_video_text_sources_video_id ON video_text_sources(video_id);
CREATE INDEX IF NOT EXISTS idx_video_text_sources_source_type ON video_text_sources(source_type);

CREATE TABLE IF NOT EXISTS recipes (
  id BIGSERIAL PRIMARY KEY,
  video_id BIGINT REFERENCES bilibili_videos(id) ON DELETE CASCADE,
  dish_name VARCHAR(255),
  servings INT,
  recipe_json JSONB NOT NULL,
  confidence NUMERIC(4, 3),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_recipes_video_id ON recipes(video_id);
CREATE INDEX IF NOT EXISTS idx_recipes_dish_name ON recipes(dish_name);

CREATE TABLE IF NOT EXISTS nutrition_foods (
  id BIGSERIAL PRIMARY KEY,
  canonical_name VARCHAR(255) NOT NULL,
  aliases JSONB DEFAULT '[]',
  kcal_per_100g NUMERIC(10, 2),
  protein_per_100g NUMERIC(10, 2),
  fat_per_100g NUMERIC(10, 2),
  carbs_per_100g NUMERIC(10, 2),
  source VARCHAR(128),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_nutrition_foods_canonical_name ON nutrition_foods(canonical_name);

CREATE TABLE IF NOT EXISTS recipe_nutrition_results (
  id BIGSERIAL PRIMARY KEY,
  recipe_id BIGINT REFERENCES recipes(id) ON DELETE CASCADE,
  total_kcal NUMERIC(10, 2),
  kcal_per_serving NUMERIC(10, 2),
  result_json JSONB NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_recipe_nutrition_results_recipe_id ON recipe_nutrition_results(recipe_id);
