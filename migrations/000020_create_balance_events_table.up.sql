CREATE TABLE IF NOT EXISTS balance_events
(
    id             BIGSERIAL                NOT NULL,
    created_at     TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at     TIMESTAMP WITH TIME ZONE NOT NULL,

    height          DECIMAL(65, 0)           NOT NULL,
    address         TEXT                     NOT NULL,
    escrow_address  TEXT                     NOT NULL,
    amount          BIGINT                   NOT NULL,
    kind            TEXT                     NOT NULL,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_balance_events_height on balance_events (height);
CREATE index idx_balance_events_address on balance_events (address);
