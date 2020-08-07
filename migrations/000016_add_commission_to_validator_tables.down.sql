ALTER TABLE validator_aggregates DROP COLUMN recent_commission;

ALTER TABLE validator_sequences DROP COLUMN commission;

ALTER TABLE validator_summary DROP COLUMN commission_avg;
ALTER TABLE validator_summary DROP COLUMN commission_min;
ALTER TABLE validator_summary DROP COLUMN commission_max;