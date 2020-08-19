ALTER TABLE validator_aggregates ADD COLUMN recent_rewards DECIMAL(65, 0);

ALTER TABLE validator_sequences ADD COLUMN rewards DECIMAL(65, 0);
