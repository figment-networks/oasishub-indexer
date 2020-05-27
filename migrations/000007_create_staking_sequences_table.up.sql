CREATE TABLE IF NOT EXISTS staking_sequences
(
    id                    BIGSERIAL                NOT NULL,
    created_at            TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at            TIMESTAMP WITH TIME ZONE NOT NULL,

    height                DECIMAL(65, 0)           NOT NULL,
    time                  TIMESTAMP WITH TIME ZONE NOT NULL,

    total_supply          DECIMAL(65, 0)           NOT NULL,
    common_pool           DECIMAL(65, 0)           NOT NULL,
    debonding_interval    BIGINT                   NOT NULL,
    min_delegation_amount DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (time, id)
);

-- Hypertable
SELECT create_hypertable('staking_sequences', 'time', if_not_exists => TRUE);

-- Indexes
CREATE index idx_staking_sequences_height on staking_sequences (height, time DESC);
