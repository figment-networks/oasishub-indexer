CREATE TABLE IF NOT EXISTS syncables
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    height        DECIMAL(65, 0)           NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,
    app_version   BIGINT                   NOT NULL,
    block_version BIGINT                   NOT NULL,

    status        SMALLINT DEFAULT 0,
    report_id     BIGINT,
    started_at    TIMESTAMP WITH TIME ZONE,
    processed_at  TIMESTAMP WITH TIME ZONE,
    duration      BIGINT,
    details       JSONB,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_syncables_report_id on syncables (report_id);
CREATE index idx_syncables_height on syncables (height);
CREATE index idx_syncables_processed_at on syncables (processed_at);