package store

const (
	summarizeBalanceQuerySelect = `
	s.time_bucket,
	s.start_height,
	balance_events.address,
	balance_events.escrow_address,
	SUM(case when kind = 'reward' then amount else 0 end) as total_rewards,
	SUM(case when kind = 'commission' then amount else 0 end) as total_commission,
	SUM(case when kind = 'slash_active' or kind = 'slash_debonding' then amount else 0 end) as total_slashed
`
	summarizeBalanceJoinQuery = `INNER JOIN
(
	SELECT
	  MAX(height)     AS end_height,
	  MIN(height)     AS start_height,
	  DATE_TRUNC(?, time) as time_bucket
	FROM syncables
	GROUP BY time_bucket
 ) AS s ON balance_events.height >= s.start_height AND balance_events.height <= s.end_height`
)
