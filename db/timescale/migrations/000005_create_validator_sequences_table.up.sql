CREATE TABLE IF NOT EXISTS validator_sequences
(
    id                  uuid                     NOT NULL DEFAULT uuid_generate_v4(),
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL,

    chain_id            TEXT                     NOT NULL,
    height              DOUBLE PRECISION         NOT NULL,
    time                TIMESTAMP WITH TIME ZONE NOT NULL,

    entity_uid          TEXT                     NOT NULL,
    node_uid            TEXT                     NOT NULL,
    consensus_uid       TEXT                     NOT NULL,
    voting_power        DOUBLE PRECISION         NOT NULL,
    total_shares        DECIMAL(65, 0)           NOT NULL,
    proposed            BOOLEAN                  NOT NULL,
    address             TEXT                     NOT NULL,
    precommit_validated BOOLEAN,
    precommit_type      NUMERIC,
    precommit_index     DOUBLE PRECISION,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('validator_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_validator_sequences_chain_id on validator_sequences (chain_id, time DESC);
CREATE index idx_validator_sequences_height on validator_sequences (height, time DESC);
CREATE index idx_validator_sequences_validator_id on validator_sequences (entity_uid, time DESC);
CREATE index idx_validator_sequences_node_uid on validator_sequences (node_uid, time DESC);
CREATE index idx_validator_sequences_proposed on validator_sequences (proposed, time DESC);
CREATE index idx_validator_sequences_total_shares on validator_sequences (total_shares, time DESC);
