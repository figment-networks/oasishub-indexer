CREATE TABLE IF NOT EXISTS reports
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    start_height  DECIMAL(65, 0)           NOT NULL,
    end_height    DECIMAL(65, 0)           NOT NULL,
    success_count INT,
    error_count   INT,
    error_msg     TEXT,
    duration      BIGINT,
    details       JSONB,
    completed_at  TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (created_at, id)
);

-- Hypertable
SELECT create_hypertable('reports', 'created_at', if_not_exists => TRUE);

-- Indexes
