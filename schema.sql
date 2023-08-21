CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY NOT NULL,
  username TEXT NOT NULL,
  hash TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS servers (
  id TEXT PRIMARY KEY NOT NULL,
  name TEXT NOT NULL,
  runtime TEXT NOT NULL,
  path TEXT NOT NULL,
  pid INTEGER DEFAULT -1,
  status BOOLEAN DEFAULT false NOT NULL,
  user_id INTEGER NOT NULL,
  FOREIGN KEY (user_id)
    REFERENCES users (id)
);
