-- +migrate Up

CREATE TABLE file_metadata (
     id UUID PRIMARY KEY,
     key VARCHAR(255) NOT NULL,
     bucket_name VARCHAR(255) NOT NULL,
     metadata JSONB NOT NULL DEFAULT '{}',
     created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_file_metadata_bucket FOREIGN KEY (bucket_name) REFERENCES buckets(name) ON DELETE CASCADE
);

CREATE INDEX idx_file_metadata_key ON file_metadata(key);
CREATE INDEX idx_file_metadata_bucket_name ON file_metadata(bucket_name);
CREATE INDEX idx_file_metadata_bucket_key ON file_metadata(bucket_name, key);

-- +migrate Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP INDEX idx_file_metadata_bucket_key;
DROP INDEX idx_file_metadata_bucket_name;
DROP INDEX idx_file_metadata_key;
DROP TABLE file_metadata;