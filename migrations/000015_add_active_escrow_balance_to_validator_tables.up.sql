ALTER TABLE validator_aggregates ADD COLUMN recent_active_escrow_balance DECIMAL(65, 0);

ALTER TABLE validator_sequences ADD COLUMN active_escrow_balance DECIMAL(65, 0);

ALTER TABLE validator_summary ADD COLUMN active_escrow_balance_avg DECIMAL(65, 0);
ALTER TABLE validator_summary ADD COLUMN active_escrow_balance_min DECIMAL(65, 0);
ALTER TABLE validator_summary ADD COLUMN active_escrow_balance_max DECIMAL(65, 0);