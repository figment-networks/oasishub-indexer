CREATE TABLE IF NOT EXISTS validator_sequences
(
    id                      BIGSERIAL                NOT NULL,

    height                  DECIMAL(65, 0)           NOT NULL,
    time                    TIMESTAMP WITH TIME ZONE NOT NULL,

    entity_uid              TEXT                     NOT NULL,
    voting_power            DECIMAL(65, 0)           NOT NULL,
    total_shares            DECIMAL(65, 0)           NOT NULL,
    proposed                BOOLEAN                  NOT NULL,
    address                 TEXT                     NOT NULL,
    precommit_validated     BOOLEAN,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('validator_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_validator_sequences_height on validator_sequences (height, time DESC);
CREATE index idx_validator_sequences_validator_id on validator_sequences (entity_uid, time DESC);
