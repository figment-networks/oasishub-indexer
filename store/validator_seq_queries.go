package store

import "fmt"

const (
	validatorSummaryForInterval = `
SELECT * 
FROM validators_summary_%s 
WHERE time_interval >= (
	SELECT time_interval 
	FROM validators_summary_%s 
	ORDER BY time_interval DESC 
	LIMIT 1
) - ?::INTERVAL
	AND entity_uid = ?
`

	allValidatorsSummaryForInterval = `
SELECT
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
FROM validators_summary_%s
WHERE time_interval >= (
	SELECT time_interval 
	FROM validators_summary_%s 
	ORDER BY time_interval DESC 
	LIMIT 1
) - ?::INTERVAL
GROUP BY time_interval
ORDER BY time_interval
`

	deleteOldValidatorSeqQuery = `
DELETE 
FROM validator_sequences 
WHERE time < NOW() - ?::INTERVAL
`

	deleteOldValidatorHourlySummary = `
DELETE 
FROM validators_summary_%s 
WHERE time < NOW() - ?::INTERVAL
`
)

func ValidatorSummaryForIntervalQuery(interval string) string {
	return fmt.Sprintf(validatorSummaryForInterval, interval, interval)
}

func AllValidatorsSummaryForIntervalQuery(interval string) string {
	return fmt.Sprintf(allValidatorsSummaryForInterval, interval, interval)
}

func DeleteOldValidatorHourlySummaryQuery(interval string) string {
	return fmt.Sprintf(deleteOldValidatorHourlySummary, interval)
}