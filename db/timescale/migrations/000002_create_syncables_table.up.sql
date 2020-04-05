CREATE TABLE IF NOT EXISTS syncables
(
    id                   BIGSERIAL NOT NULL,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at           TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id             TEXT NOT NULL,
    height               DOUBLE PRECISION NOT NULL,
    time                 TIMESTAMP WITH TIME ZONE NOT NULL,

    report_id            NUMERIC,
    type                 VARCHAR(100),
    data                 JSONB,
    processed_at         TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('syncables', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_syncables_report_id on syncables (report_id, time DESC);
CREATE index idx_syncables_chain_id on syncables (chain_id, time DESC);
CREATE index idx_syncables_height on syncables (height, time DESC);
CREATE index idx_syncables_processed_at on syncables (time DESC, processed_at);