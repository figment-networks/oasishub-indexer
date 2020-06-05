-- Hourly view
CREATE VIEW blocks_summary_hourly WITH (timescaledb.continuous) AS
SELECT
    time_bucket(INTERVAL '1 hour', time) AS time_interval,
    COUNT(*) AS count,
    EXTRACT(EPOCH FROM (last(time, time) - first(time, time)) / COUNT(*)) AS avg
FROM block_sequences
GROUP BY time_interval;

-- Daily view
CREATE VIEW blocks_summary_daily WITH (timescaledb.continuous) AS
SELECT
    time_bucket(INTERVAL '1 day', time) AS time_interval,
    COUNT(*) AS count,
    EXTRACT(EPOCH FROM (last(time, time) - first(time, time)) / COUNT(*)) AS avg
FROM block_sequences
GROUP BY time_interval;

-- Only on enterprise edition
-- SELECT add_drop_chunks_policy('block_sequences', INTERVAL '26 hours');