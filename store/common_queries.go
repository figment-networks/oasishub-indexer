package store

import "fmt"

const (
	activityPeriodsQuery = `
WITH cte AS (
    SELECT
      time_bucket,
      sum(CASE WHEN diff IS NULL OR diff > ? :: INTERVAL
        THEN 1
          ELSE NULL END)
      OVER (
        ORDER BY time_bucket ) AS period
    FROM (
           SELECT
             time_bucket,
             time_bucket - lag(time_bucket, 1)
             OVER (
               ORDER BY time_bucket ) AS diff
           FROM %v
           WHERE time_interval = ? AND index_version = ?
         ) AS x
)
SELECT
  period,
  MIN(time_bucket),
  MAX(time_bucket)
FROM cte
GROUP BY period
ORDER BY period
`
)

func getActivityPeriodsQuery(tableName string) string {
	return fmt.Sprintf(activityPeriodsQuery, tableName)
}
