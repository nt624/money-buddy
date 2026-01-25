CREATE TABLE users (
  id TEXT PRIMARY KEY,          -- Firebase UID
  income INT NOT NULL,           -- 月収（手取り）
  saving_goal INT NOT NULL,      -- 月の貯金額
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
