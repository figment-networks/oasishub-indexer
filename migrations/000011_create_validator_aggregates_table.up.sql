CREATE TABLE IF NOT EXISTS validator_aggregates
(
    id                         BIGSERIAL                NOT NULL,
    created_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at                 TIMESTAMP WITH TIME ZONE NOT NULL,

    started_at_height          DECIMAL(65, 0)           NOT NULL,
    started_at                 TIMESTAMP WITH TIME ZONE NOT NULL,
    recent_at_height           DECIMAL(65, 0)           NOT NULL,
    recent_at                  TIMESTAMP WITH TIME ZONE NOT NULL,

    entity_uid                 TEXT,
    recent_address             TEXT,
    recent_voting_power        DECIMAL(65, 0),
    recent_total_shares        DECIMAL(65, 0),
    recent_as_validator_height DECIMAL(65, 0),
    recent_proposed_height     DECIMAL(65, 0),
    accumulated_proposed_count BIGINT,
    accumulated_uptime         BIGINT,
    accumulated_uptime_count   BIGINT,

    PRIMARY KEY (id)
);

-- Indexes
CREATE index idx_validator_aggregates_entity_uid on validator_aggregates (entity_uid);