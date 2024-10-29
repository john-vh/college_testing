CREATE TYPE post_status AS ENUM ('active', 'disabled', 'archived');

CREATE TABLE IF NOT EXISTS posts (
  business_id UUID NOT NULL,
  id SERIAL NOT NULL,
  title VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  status post_status NOT NULL DEFAULT 'disabled',
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(business_id, id),
  FOREIGN KEY(business_id) REFERENCES businesses(id)
);
