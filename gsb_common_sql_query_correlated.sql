SELECT
  id,
  title,
  (SELECT COUNT(*) FROM comments WHERE post_id = posts.id) as comment_count,
  (SELECT AVG(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS avg_comment_length,
  (SELECT MAX(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS max_comment_length
FROM posts
