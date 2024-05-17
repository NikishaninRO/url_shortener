BEGIN;

CREATE TABLE
  urls (
    id SERIAL PRIMARY KEY,
    alias TEXT NOT NULL UNIQUE,
    url TEXT NOT NULL
  );

CREATE INDEX idx_alias ON urls (alias);

COMMIT;