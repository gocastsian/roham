-- +migrate Up
CREATE TABLE IF NOT EXISTS job (
    id SERIAL PRIMARY KEY,
    file_key VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    workflow_id VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    error_message TEXT,
    started_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
    );

-- +migrate StatementBegin
CREATE OR REPLACE FUNCTION update_job_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER trigger_update_job_timestamp
    BEFORE UPDATE ON job
    FOR EACH ROW
    EXECUTE FUNCTION update_job_timestamp();
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TRIGGER IF EXISTS trigger_update_job_timestamp ON job;
DROP FUNCTION IF EXISTS update_job_timestamp;
-- +migrate StatementEnd

-- +migrate Down
DROP TABLE IF EXISTS job;
