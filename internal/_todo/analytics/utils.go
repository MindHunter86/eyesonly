package analytics

const sqlDatabaseDriver string = "duckdb"
const sqlDatabaseParams string = "?access_mode=READ_WRITE"

// const dbAnalyzeMacroTpl string = ""

var dbSchemaMacros []string = []string{
	`CREATE OR REPLACE MACRO get_per_second_values(tbl) AS TABLE
SELECT
  time,
  FLOOR(SUM(delta) OVER (
    ORDER BY time
    RANGE BETWEEN INTERVAL 1 MINUTE PRECEDING AND CURRENT ROW
  ) / 60.0) AS rps_60s
FROM (
  SELECT
    time,
    GREATEST(0, value - LAG(value, 1, value) OVER (ORDER BY time)) AS delta
  FROM query_table(tbl)
) d`,

	`CREATE OR REPLACE MACRO get_1m_shifted_differences(with_time, metric_name) AS TABLE
WITH src AS (
  SELECT time, value
  FROM stats
  WHERE time BETWEEN with_time - INTERVAL 2 MINUTES AND now()
    AND metric = metric_name
),
curr AS (
  SELECT * FROM get_per_second_values(src)
),
shift AS (
  SELECT time + INTERVAL 1 MINUTE AS time, rps_60s
  FROM curr
)
SELECT
  curr.time AS time,
  curr.rps_60s AS rps_now,
  COALESCE(shift.rps_60s, 0) AS rps_shifted,
  COALESCE(curr.rps_60s - shift.rps_60s, 0) AS diff,
  CASE WHEN COALESCE(shift.rps_60s, 0) > 0
    THEN ROUND(curr.rps_60s / shift.rps_60s, 3)
    ELSE 0.0
  END AS ratio
FROM curr
ASOF LEFT JOIN shift
  ON shift.time <= curr.time
WHERE curr.time BETWEEN with_time AND now()
ORDER BY curr.time DESC
LIMIT 4`,

	`CREATE OR REPLACE MACRO get_1m_requested_differences(diff_time, from_time, metric_name) AS TABLE
WITH src AS (
  SELECT time, value
  FROM stats
  WHERE time BETWEEN diff_time - INTERVAL 1 MINUTE AND COALESCE(from_time, NOW())
    AND metric = metric_name
),
rps AS (
  SELECT * FROM get_per_second_values(src)
  WHERE time BETWEEN diff_time AND COALESCE(from_time, NOW())
),
agg AS (
  SELECT
    MAX(time) as time_now,
    MIN(time) as time_prev,
    ARG_MAX(rps_60s, time) as rps_now,
    ARG_MIN(rps_60s, time) as rps_prev
  FROM rps
)
SELECT
  time_now,
  time_prev,
  COALESCE(rps_now, 0) as rps_now,
  COALESCE(rps_prev, 0) as rps_prev,
  COALESCE(rps_now - rps_prev, 0) as rps_diff,
  CASE WHEN rps_prev > 0
    THEN ROUND(rps_now / rps_prev, 3)
    ELSE 0.0
  END AS ratio
FROM agg`,
}

//

// SELECT time, FLOOR(SUM(delta) OVER ( ORDER BY time RANGE BETWEEN INTERVAL 60 SECONDS PRECEDING AND CURRENT ROW ) / 60.0) AS rps_60s
// FROM ( SELECT time, GREATEST( 0, value - LAG(value, 1, value) OVER (ORDER BY time) ) AS delta
//   FROM ( SELECT time, MAX(value) AS value FROM stats WHERE metric = 'fiber.requests' GROUP BY time ) s ) t ORDER BY time;

/*

CREATE OR REPLACE MACRO rps60_series(tbl) AS TABLE (
  SELECT
    metric,
    time,
    FLOOR(SUM(delta) OVER (
      PARTITION BY metric
      ORDER BY time
      RANGE BETWEEN INTERVAL 60 SECONDS PRECEDING AND CURRENT ROW
    ) / 60.0) AS rps_60s
  FROM (
    SELECT
      metric,
      time,
      GREATEST(
        0,
        value - LAG(value) OVER (PARTITION BY metric ORDER BY time)
      ) AS delta
    FROM tbl
  ) d
);*/
