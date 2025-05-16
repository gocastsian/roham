-- +migrate Up

CREATE TABLE file_metadata
(
    id         SERIAL PRIMARY KEY,
    file_key   VARCHAR(255) NOT NULL UNIQUE,
    storage_id BIGINT       NOT NULL REFERENCES storages (id) ON DELETE RESTRICT,
    file_name  VARCHAR(255) NOT NULL,
    mime_type  VARCHAR(255) NOT NULL,
    file_size  BIGINT       NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT now(),
    claimed_at TIMESTAMPTZ
);

CREATE INDEX idx_file_metadata_storage_id ON file_metadata (storage_id);
CREATE INDEX idx_file_metadata_created_at ON file_metadata (created_at);
CREATE INDEX idx_file_metadata_claimed_at ON file_metadata (claimed_at);

-- +migrate Down

DROP INDEX IF EXISTS idx_file_metadata_claimed_at;
DROP INDEX IF EXISTS idx_file_metadata_created_at;
DROP INDEX IF EXISTS idx_file_metadata_storage_id;

DROP TABLE IF EXISTS file_metadata;