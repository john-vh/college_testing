CREATE TYPE post_application_status AS ENUM ('pending', 'accepted', 'rejected', 'withdrawn', 'completed');

CREATE TABLE IF NOT EXISTS post_applications (
  business_id UUID NOT NULL,
  post_id INT NOT NULL,
  user_id UUID NOT NULL,
  status post_application_status NOT NULL DEFAULT 'pending',
  notes TEXT NOT NULL DEFAULT '',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(business_id, post_id, user_id),
  FOREIGN KEY(business_id, post_id) REFERENCES posts(business_id, id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
