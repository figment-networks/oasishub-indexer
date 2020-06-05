-- Hourly views
CREATE VIEW validators_summary_hourly WITH (timescaledb.continuous) AS
SELECT
    entity_uid,
    time_bucket(INTERVAL '1 hour', time) AS time_interval,
    AVG(voting_power) AS voting_power_avg,
    MAX(voting_power) AS voting_power_max,
    MIN(voting_power) AS voting_power_min,
    AVG(total_shares) AS total_shares_avg,
    MAX(total_shares) AS total_shares_max,
    MIN(total_shares) AS total_shares_min,
    AVG(precommit_validated::INT) AS uptime_avg,
    SUM(precommit_validated::INT) AS validated_sum,
    COUNT(*) - SUM(precommit_validated::INT) AS not_validated_sum,
    SUM(proposed::INT) AS proposed_sum
FROM validator_sequences
GROUP BY entity_uid, time_interval;

-- Daily view
CREATE VIEW validators_summary_daily WITH (timescaledb.continuous) AS
SELECT
    entity_uid,
    time_bucket(INTERVAL '1 day', time) AS time_interval,
    AVG(voting_power) AS voting_power_avg,
    MAX(voting_power) AS voting_power_max,
    MIN(voting_power) AS voting_power_min,
    AVG(total_shares) AS total_shares_avg,
    MAX(total_shares) AS total_shares_max,
    MIN(total_shares) AS total_shares_min,
    AVG(precommit_validated::INT) AS uptime_avg,
    SUM(precommit_validated::INT) AS validated_sum,
    COUNT(*) - SUM(precommit_validated::INT) AS not_validated_sum,
    SUM(proposed::INT) AS proposed_sum
FROM validator_sequences
GROUP BY entity_uid, time_interval;

-- Only on enterprise edition
-- SELECT add_drop_chunks_policy('validator_sequences', INTERVAL '26 hours');