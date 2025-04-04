SELECT
  date(created) as day,
  SUM(json_extract(stats, '$.lorem')) as sum_lorem,
  AVG(json_extract(stats, '$.ipsum.dolor')) as avg_dolor,
  MAX(json_extract(stats, '$.lorem.sit')) as max_sit
FROM posts
GROUP BY day
ORDER BY day
