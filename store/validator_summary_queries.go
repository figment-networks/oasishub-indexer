package store

const (
	validatorSummaryForIntervalQuery = `
SELECT * 
FROM validator_summary 
WHERE time_bucket >= (
	SELECT time_bucket 
	FROM validator_summary 
	WHERE time_interval = ?
	ORDER BY time_bucket DESC
	LIMIT 1
) - ?::INTERVAL
	AND address = ? AND time_interval = ?
ORDER BY time_bucket
`

	allValidatorsSummaryForIntervalQuery = `
SELECT
  time_bucket,
  time_interval,
  AVG(voting_power_avg) AS voting_power_avg,
  MAX(voting_power_max) AS voting_power_max,
  MIN(voting_power_min) AS voting_power_min,
  AVG(total_shares_avg) AS total_shares_avg,
  MAX(total_shares_max) AS total_shares_max,
  MIN(total_shares_min) AS total_shares_min,
  AVG(active_escrow_balance_avg) AS active_escrow_balance_avg,
  MAX(active_escrow_balance_max) AS active_escrow_balance_max,
  MIN(active_escrow_balance_min) AS active_escrow_balance_min,
  AVG(commission_avg) AS commission_avg,
  MAX(commission_max) AS commission_max,
  MIN(commission_min) AS commission_min,
  AVG(uptime_avg) AS uptime_avg,
  SUM(validated_sum) AS validated_sum,
  SUM(not_validated_sum) AS not_validated_sum,
  SUM(proposed_sum) AS proposed_sum
FROM validator_summary
WHERE time_bucket >= (
	SELECT time_bucket 
	FROM validator_summary 
	WHERE time_interval = ?
	ORDER BY time_bucket DESC 
	LIMIT 1
) - ?::INTERVAL
	AND time_interval = ?
GROUP BY time_bucket, time_interval
ORDER BY time_bucket
`
)
