package store

const (
	summarizeBalanceQuerySelect = `
	s.time_bucket,
	s.min_height,
	balance_events.kind,
	balance_events.address,
	balance_events.escrow_address,
	SUM(balance_events.amount) as total_amount
`
	summarizeBalanceJoinQuery = `INNER JOIN
(
	SELECT
	  MAX(height)     AS max_height,
	  MIN(height)     AS min_height,
	  DATE_TRUNC(?, time) as time_bucket
	FROM syncables
	GROUP BY time_bucket
 ) AS s ON balance_events.height >= s.min_height AND balance_events.height <= s.max_height`
)
