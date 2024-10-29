CREATE TYPE business_status AS ENUM ('pending', 'active', 'disabled');

CREATE TABLE IF NOT EXISTS businesses (
  user_id UUID NOT NULL,
  id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  website VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  status business_status NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(id),
  FOREIGN KEY(user_id) REFERENCES users(id),
  UNIQUE(name)
);
