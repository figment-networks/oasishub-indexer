CREATE TABLE IF NOT EXISTS block_sequences
(
    id                  BIGSERIAL                NOT NULL,

    height              DECIMAL(65, 0)           NOT NULL,
    time                TIMESTAMP WITH TIME ZONE NOT NULL,

    transactions_count  INT,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_block_sequences_height on block_sequences (height);
