package store


const (
	totalSharesForIntervalQuery = `
SELECT
  time_bucket($1, time) AS time_interval,
  SUM(a) as sum,
  COUNT(*) as count,
  SUM(a) / COUNT(*) AS avg
FROM (
  SELECT
    MAX(time) as time,
    SUM(total_shares) / COUNT(*) AS a
  FROM validator_sequences
    WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $2::INTERVAL
  GROUP BY height
  ORDER BY height
) d
GROUP BY time_interval
ORDER BY time_interval;
`
	totalVotingPowerForIntervalQuery = `
SELECT
  time_bucket($1, time) AS time_interval,
  SUM(a) as sum,
  COUNT(*) as count,
  SUM(a) / COUNT(*) AS avg
FROM (
  SELECT
    MAX(time) as time,
    SUM(voting_power) / COUNT(*) AS a
  FROM validator_sequences
    WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $2::INTERVAL
  GROUP BY height
  ORDER BY height
) d
GROUP BY time_interval
ORDER BY time_interval;
`
	validatorSharesForIntervalQuery=`
SELECT
  time_bucket($2, time) AS time_interval,
  AVG(total_shares) AS avg
FROM validator_sequences
  WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $3::INTERVAL AND entity_uid = $1
GROUP BY time_interval
ORDER BY time_interval ASC;
`
	validatorVotingPowerForIntervalQuery=`
SELECT
  time_bucket($2, time) AS time_interval,
  AVG(voting_power) AS avg
FROM validator_sequences
  WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $3::INTERVAL AND entity_uid = $1
GROUP BY time_interval
ORDER BY time_interval ASC;
`
	validatorUptimeForIntervalQuery=`
SELECT
  time_bucket($2, time) AS time_interval,
  AVG(precommit_validated::INT) AS avg
FROM validator_sequences
  WHERE (
      SELECT time
      FROM validator_sequences
      ORDER BY time DESC
      LIMIT 1
    ) - time < $3::INTERVAL AND entity_uid = $1
GROUP BY time_interval
ORDER BY time_interval ASC;
`
)
