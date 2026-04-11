-- Migration: 002_monthly_settings
-- Purpose: 手取り・目標額を月毎に設定できるようにする
-- ※ このファイルを実行する前に必ずバックアップを取ること

-- ============================================================
-- Step 1: monthly_settings テーブル作成
-- ============================================================
CREATE TABLE IF NOT EXISTS monthly_settings (
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
