WITH RECURSIVE dates(day) AS (
  SELECT date('now', '-30 days')
  UNION ALL
  SELECT date(day, '+1 day')
  FROM dates
  WHERE day < date('now')
)
SELECT day,
  (SELECT COUNT(*) FROM posts WHERE date(created) = day) as day_total
FROM dates
ORDER BY day
