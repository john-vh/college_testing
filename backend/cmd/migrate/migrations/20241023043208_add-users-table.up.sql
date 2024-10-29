CREATE TYPE user_status AS ENUM ('active', 'banned', 'disabled');

CREATE TABLE IF NOT EXISTS users (
  id UUID NOT NULL,
  status user_status NOT NULL DEFAULT 'active',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(id)
);
