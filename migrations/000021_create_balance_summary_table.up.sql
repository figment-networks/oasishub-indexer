CREATE TABLE IF NOT EXISTS balance_summary
(
    id             BIGSERIAL                NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL,

    time_interval  VARCHAR                  NOT NULL,
    time_bucket    TIMESTAMP WITH TIME ZONE NOT NULL,
    index_version  INT                      NOT NULL,

    start_height     DECIMAL(65, 0)         NOT NULL,
    address          TEXT                   NOT NULL,
    escrow_address   TEXT                   NOT NULL,
    total_rewards    BIGINT,
    total_commission BIGINT,
    total_slashed    BIGINT,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_balance_summary_time on balance_summary (time_interval, time_bucket);
CREATE index idx_balance_summary_index_version on balance_summary (index_version);
CREATE index idx_balance_summary_address on balance_summary (address);

