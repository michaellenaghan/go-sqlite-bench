SELECT
  strftime('%Y-%m', created) AS month,
  COUNT(*) as month_total
FROM posts
GROUP BY month
ORDER BY month
