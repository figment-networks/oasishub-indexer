CREATE TABLE IF NOT EXISTS entity_aggregates
(
--     Domain entity
    id                BIGSERIAL                NOT NULL,
    created_at        TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at        TIMESTAMP WITH TIME ZONE NOT NULL,

-- Aggregate
    started_at_height NUMERIC                  NOT NULL,
    started_at        TIMESTAMP WITH TIME ZONE NOT NULL,

-- Custom
    entity_uid        TEXT,

    PRIMARY KEY (id)
);

-- Hypertable

-- Indexes
CREATE index idx_entity_aggregates_entity_uid on entity_aggregates (entity_uid);