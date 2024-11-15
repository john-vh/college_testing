CREATE TYPE user_role AS ENUM ('admin', 'user');

CREATE TABLE IF NOT EXISTS user_roles(
  user_id UUID NOT NULL,
  role user_role NOT NULL,

  PRIMARY KEY(user_id, role)
);
