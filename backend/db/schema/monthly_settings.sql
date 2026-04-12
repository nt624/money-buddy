CREATE TABLE monthly_settings (
  id          SERIAL PRIMARY KEY,
  user_id     TEXT    NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  year        INT     NOT NULL,
  month       INT     NOT NULL,
  income      INT     NOT NULL,
  saving_goal INT     NOT NULL,
  created_at  TIMESTAMP NOT NULL DEFAULT now(),
  updated_at  TIMESTAMP NOT NULL DEFAULT now(),
  CONSTRAINT monthly_settings_month_check CHECK (month BETWEEN 1 AND 12),
  CONSTRAINT monthly_settings_user_year_month_unique UNIQUE (user_id, year, month)
);

CREATE INDEX IF NOT EXISTS idx_monthly_settings_user_id ON monthly_settings(user_id);
