package store

const (
	summarizeValidatorsQuerySelect = `
	entity_uid,
	DATE_TRUNC(?, time)                      AS time_bucket,
   	AVG(voting_power)                        AS voting_power_avg,
   	MAX(voting_power)                        AS voting_power_max,
   	MIN(voting_power)                        AS voting_power_min,
   	AVG(total_shares)                        AS total_shares_avg,
   	MAX(total_shares)                        AS total_shares_max,
   	MIN(total_shares)                        AS total_shares_min,
   	AVG(precommit_validated::INT)            AS uptime_avg,
   	SUM(precommit_validated::INT)            AS validated_sum,
   	COUNT(*) - SUM(precommit_validated::INT) AS not_validated_sum,
   	SUM(proposed::INT)                       AS proposed_sum
`
)
