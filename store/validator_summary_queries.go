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
	AND entity_uid = ? AND time_interval = ?
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
