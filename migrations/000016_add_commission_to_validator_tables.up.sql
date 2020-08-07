ALTER TABLE validator_aggregates ADD COLUMN recent_commission DECIMAL(65, 0);

ALTER TABLE validator_sequences ADD COLUMN commission DECIMAL(65, 0);

ALTER TABLE validator_summary ADD COLUMN commission_avg DECIMAL(65, 0);
ALTER TABLE validator_summary ADD COLUMN commission_min DECIMAL(65, 0);
ALTER TABLE validator_summary ADD COLUMN commission_max DECIMAL(65, 0);