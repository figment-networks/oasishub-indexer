CREATE TABLE IF NOT EXISTS block_sequences
(
    id                  BIGSERIAL                NOT NULL,

    height              DECIMAL(65, 0)           NOT NULL,
    time                TIMESTAMP WITH TIME ZONE NOT NULL,

    transactions_count  INT,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('block_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_block_sequences_height on block_sequences (height, time DESC);
