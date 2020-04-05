CREATE TABLE IF NOT EXISTS delegation_sequences
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id      TEXT                     NOT NULL,
    height        DOUBLE PRECISION         NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,

    validator_uid TEXT                     NOT NULL,
    delegator_uid TEXT                     NOT NULL,
    shares        DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('delegation_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_delegation_sequences_chain_id on delegation_sequences (chain_id, time DESC);
CREATE index idx_delegation_sequences_height on delegation_sequences (height, time DESC);
CREATE index idx_delegation_sequences_app_version on delegation_sequences (validator_uid, time DESC);
CREATE index idx_delegation_sequences_block_version on delegation_sequences (delegator_uid, time DESC);
