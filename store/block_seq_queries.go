package store

import "fmt"

const (
	blockTimesForRecentBlocksQuery = `
SELECT 
  MIN(height) start_height, 
  MAX(height) end_height, 
  MIN(time) start_time,
  MAX(time) end_time,
  COUNT(*) count, 
  EXTRACT(EPOCH FROM MAX(time) - MIN(time)) AS diff, 
  EXTRACT(EPOCH FROM ((MAX(time) - MIN(time)) / COUNT(*))) AS avg
  FROM ( 
    SELECT * FROM block_sequences
    ORDER BY height DESC
    LIMIT ?
  ) t;
`

	allBlocksSummaryForInterval = `
SELECT * 
FROM blocks_summary_%s 
WHERE time_interval >= (
	SELECT time_interval 
	FROM blocks_summary_%s 
	ORDER BY time_interval DESC 
	LIMIT 1
) - ?::INTERVAL
`

	deleteOldBlockSeqQuery = `
DELETE 
FROM block_sequences 
WHERE time < NOW() - ?::INTERVAL
`

	deleteOldBlockHourlySummary = `
DELETE 
FROM blocks_summary_%s 
WHERE time < NOW() - ?::INTERVAL
`
)

func AllBlocksSummaryForIntervalQuery(interval string) string {
	return fmt.Sprintf(allBlocksSummaryForInterval, interval, interval)
}

func DeleteOldBlockHourlySummaryQuery(interval string) string {
	return fmt.Sprintf(deleteOldBlockHourlySummary, interval)
}