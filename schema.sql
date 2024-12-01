CREATE TABLE secret (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT UNIQUE,
  value TEXT,
  created_at INTEGER,
  owner TEXT
);

CREATE INDEX IF NOT EXISTS name_idx ON secret (name);
