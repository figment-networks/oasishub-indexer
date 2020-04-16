CREATE TABLE IF NOT EXISTS transaction_sequences
(
    id         BIGSERIAL                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id   TEXT                     NOT NULL,
    height     DOUBLE PRECISION         NOT NULL,
    time       TIMESTAMP WITH TIME ZONE NOT NULL,

    public_key TEXT                     NOT NULL,
    hash       TEXT                     NOT NULL,
    nonce      NUMERIC                  NOT NULL,
    fee        DECIMAL(65, 0)           NOT NULL,
    gas_limit  DOUBLE PRECISION         NOT NULL,
    gas_price  DECIMAL(65, 0)           NOT NULL,
    method     TEXT                     NOT NULL,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('transaction_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_transaction_sequences_chain_id on transaction_sequences (chain_id, time DESC);
CREATE index idx_transaction_sequences_height on transaction_sequences (height, time DESC);
CREATE index idx_transaction_sequences_public_key on transaction_sequences (public_key, time DESC);
CREATE index idx_transaction_sequences_hash on transaction_sequences (hash, time DESC);
