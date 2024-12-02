-- 1/12/2024
CREATE TABLE IF NOT EXISTS secret (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE,
  value TEXT,
  created_at INTEGER,
  owner TEXT
);

CREATE INDEX IF NOT EXISTS name_idx ON secret (name);

-- 2/12/2024
-- add encryption related columns, only ever store ciphered value, iv and salt
-- even non passphrase values are encrypted
ALTER TABLE secret ADD iv TEXT;
ALTER TABLE secret ADD salt TEXT;
ALTER TABLE secret ADD has_password NUMBER;

-- removed unused column
ALTER TABLE secret DROP owner;
