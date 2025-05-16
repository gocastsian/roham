
-- +migrate Up
CREATE TABLE buckets (
  id UUID PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE
);

CREATE INDEX idx_buckets_name ON buckets(name);

-- +migrate Down
DROP INDEX idx_buckets_name;
DROP TABLE buckets;