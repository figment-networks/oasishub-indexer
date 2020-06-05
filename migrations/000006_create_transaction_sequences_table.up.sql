CREATE TABLE IF NOT EXISTS transaction_sequences
(
    id         BIGSERIAL                NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL,

    height     DECIMAL(65, 0)           NOT NULL,
    time       TIMESTAMP WITH TIME ZONE NOT NULL,

    public_key TEXT                     NOT NULL,
    hash       TEXT                     NOT NULL,
    nonce      BIGINT                   NOT NULL,
    fee        DECIMAL(65, 0)           NOT NULL,
    gas_limit  DECIMAL(65, 0)           NOT NULL,
    gas_price  DECIMAL(65, 0)           NOT NULL,
    method     TEXT                     NOT NULL,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_transaction_sequences_height on transaction_sequences (height);
CREATE index idx_transaction_sequences_public_key on transaction_sequences (public_key);
