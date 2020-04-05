CREATE TABLE IF NOT EXISTS reports
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    start_height  NUMERIC                  NOT NULL,
    end_height    NUMERIC                  NOT NULL,
    success_count NUMERIC,
    error_count   NUMERIC,
    error_msg     TEXT,
    duration      NUMERIC,
    details       JSONB,
    completed_at  TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (created_at, id)
);

-- Hypertable
SELECT create_hypertable('reports', 'created_at', if_not_exists => TRUE);

-- Indexes
