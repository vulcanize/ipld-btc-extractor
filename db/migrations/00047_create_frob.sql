-- +goose Up
CREATE TABLE maker.frob (
  id        SERIAL PRIMARY KEY,
  header_id INTEGER NOT NULL REFERENCES headers (id) ON DELETE CASCADE,
  ilk       TEXT,
  urn       TEXT,
  dink      NUMERIC,
  dart      NUMERIC,
  ink       NUMERIC,
  art       NUMERIC,
  iart      NUMERIC,
  log_idx   INTEGER NOT NUll,
  tx_idx    INTEGER NOT NUll,
  raw_log   JSONB,
  UNIQUE (header_id, tx_idx, log_idx)
);

ALTER TABLE public.checked_headers
  ADD COLUMN frob_checked BOOLEAN NOT NULL DEFAULT FALSE;


-- +goose Down
DROP TABLE maker.frob;

ALTER TABLE public.checked_headers
  DROP COLUMN frob_checked;