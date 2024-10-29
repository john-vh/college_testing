CREATE TABLE IF NOT EXISTS user_accounts (
  user_id UUID NOT NULL,
  account_provider VARCHAR(255) NOT NULL,
  account_id VARCHAR(255) NOT NULL,
  is_primary BOOLEAN,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(user_id, account_provider, account_id),
  FOREIGN KEY(account_provider, account_id) REFERENCES accounts(provider, id),
  FOREIGN KEY(user_id) REFERENCES users(id)
);
