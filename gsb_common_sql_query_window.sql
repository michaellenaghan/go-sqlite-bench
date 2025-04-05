WITH day_totals AS (
  SELECT date(created) as day, COUNT(*) as day_total
  FROM posts
  GROUP BY day
)
SELECT day, day_total,
  SUM(day_total) OVER (ORDER BY day) as running_total
FROM day_totals
ORDER BY day
