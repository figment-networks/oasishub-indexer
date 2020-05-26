CREATE TABLE IF NOT EXISTS account_aggregates
(
    id                                   BIGSERIAL                NOT NULL,
    created_at                           TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                           TIMESTAMP WITH TIME ZONE NOT NULL,

    started_at_height                    DECIMAL(65, 0)           NOT NULL,
    started_at                           TIMESTAMP WITH TIME ZONE NOT NULL,
    recent_at_height                     DECIMAL(65, 0)           NOT NULL,
    recent_at                            TIMESTAMP WITH TIME ZONE NOT NULL,

    public_key                           TEXT,
    recent_general_nonce                 BIGINT                   NOT NULL,
    recent_general_balance               DECIMAL(65, 0)           NOT NULL,
    recent_escrow_active_balance         DECIMAL(65, 0)           NOT NULL,
    recent_escrow_active_total_shares    DECIMAL(65, 0)           NOT NULL,
    recent_escrow_debonding_balance      DECIMAL(65, 0)           NOT NULL,
    recent_escrow_debonding_total_shares DECIMAL(65, 0)           NOT NULL,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_account_aggregates_public_key on account_aggregates (public_key);