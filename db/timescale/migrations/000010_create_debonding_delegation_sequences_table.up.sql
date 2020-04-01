CREATE TABLE IF NOT EXISTS debonding_delegation_sequences
(
    id            uuid                     NOT NULL DEFAULT uuid_generate_v4(),
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id      TEXT                     NOT NULL,
    height        DOUBLE PRECISION         NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,

    validator_uid TEXT                     NOT NULL,
    delegator_uid TEXT                     NOT NULL,
    shares        DECIMAL(65, 0)           NOT NULL,
    debond_end    NUMERIC                  NOT NULL,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('debonding_delegation_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index debonding_idx_delegation_sequences_chain_id on debonding_delegation_sequences (chain_id, time DESC);
CREATE index debonding_idx_delegation_sequences_height on debonding_delegation_sequences (height, time DESC);
CREATE index debonding_idx_delegation_sequences_app_version on debonding_delegation_sequences (validator_uid, time DESC);
CREATE index debonding_idx_delegation_sequences_block_version on debonding_delegation_sequences (delegator_uid, time DESC);
