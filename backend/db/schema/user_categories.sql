CREATE TABLE user_categories (
  id SERIAL PRIMARY KEY,
  user_id TEXT NOT NULL REFERENCES users(id),
  name TEXT NOT NULL,
  category_type TEXT NOT NULL DEFAULT 'user',
  sort_order INT NOT NULL DEFAULT 0,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT user_categories_type_check CHECK (category_type IN ('default', 'user', 'other')),
  CONSTRAINT user_categories_user_name_unique UNIQUE (user_id, name)
);

CREATE INDEX idx_user_categories_user_id ON user_categories(user_id);
