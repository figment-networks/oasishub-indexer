CREATE TABLE IF NOT EXISTS delegation_sequences
(
    id            BIGSERIAL                NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL,

    height        DECIMAL(65, 0)           NOT NULL,
    time          TIMESTAMP WITH TIME ZONE NOT NULL,

    validator_uid TEXT                     NOT NULL,
    delegator_uid TEXT                     NOT NULL,
    shares        DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_delegation_sequences_height on delegation_sequences (height);
