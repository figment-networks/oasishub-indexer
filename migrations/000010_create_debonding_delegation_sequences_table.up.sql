CREATE TABLE IF NOT EXISTS debonding_delegation_sequences
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    height        DECIMAL(65, 0)           NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,

    validator_uid TEXT                     NOT NULL,
    delegator_uid TEXT                     NOT NULL,
    shares        DECIMAL(65, 0)           NOT NULL,
    debond_end    BIGINT                   NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index debonding_idx_delegation_sequences_height on debonding_delegation_sequences (height);
