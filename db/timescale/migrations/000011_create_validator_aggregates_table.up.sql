CREATE TABLE IF NOT EXISTS validator_aggregates
(
--     Domain entity
    id                         BIGSERIAL                NOT NULL,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL,

-- Aggregate
    started_at_height          NUMERIC                  NOT NULL,
    started_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    recent_at_height           NUMERIC                  NOT NULL,
    recent_at                  TIMESTAMP WITH TIME ZONE NOT NULL,

-- Custom
    entity_uid                 TEXT,
    recent_address             TEXT,
    recent_voting_power        NUMERIC,
    recent_total_shares        DECIMAL(65, 0),
    recent_as_validator_height NUMERIC,
    recent_proposed_height     NUMERIC,
    accumulated_proposed_count NUMERIC,
    accumulated_uptime         NUMERIC,
    accumulated_uptime_count   NUMERIC,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_validator_aggregates_entity_uid on validator_aggregates (entity_uid);