ALTER TABLE validator_aggregates DROP COLUMN recent_active_escrow_balance;

ALTER TABLE validator_sequences DROP COLUMN active_escrow_balance;

ALTER TABLE validator_summary DROP COLUMN active_escrow_balance_avg;
ALTER TABLE validator_summary DROP COLUMN active_escrow_balance_min;
ALTER TABLE validator_summary DROP COLUMN active_escrow_balance_max;