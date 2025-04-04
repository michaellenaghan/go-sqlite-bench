package go_sqlite_bench

const LoremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eget sapien accumsan, commodo ligula iaculis, pretium ligula. In ac lobortis nulla. Donec lobortis metus sed mauris iaculis euismod. Ut vehicula velit vitae dolor maximus euismod. Nulla vel risus eros. Vivamus porttitor odio eleifend, imperdiet tellus sed, feugiat augue. Donec faucibus eget nunc facilisis gravida. Fusce posuere ac lacus eu rutrum. Vivamus nec nibh sed nisl maximus varius. Pellentesque in placerat eros. Vivamus efficitur in dolor nec eleifend. Proin quis nibh quis enim rutrum posuere. Aliquam odio metus, scelerisque quis massa imperdiet, tincidunt placerat orci. Maecenas posuere, ex vitae porttitor tristique, risus erat bibendum augue, nec tempus eros dolor vitae ipsum.`
const LoremIpsumJSON = `{"lorem": 10, "ipsum": {"dolor": 100, "sit": 1000}}`

var SQLForSchema = []string{
	`
	create table posts (
		id integer primary key,
		title text not null,
		content text not null,
		created text not null default (strftime('%Y-%m-%dT%H:%M:%SZ')),
		stats text not null
	);
	`,
	`
	create index post_title_idx on posts (title);
	`,
	`
	create index post_created_idx on posts (created);
	`,
	`
	create table comments (
		id integer primary key,
		post_id int not null references posts (id)
			on update restrict
			on delete cascade,
		name text not null,
		content text not null,
		created text not null default (strftime('%Y-%m-%dT%H:%M:%SZ')),
		stats text not null
	);
	`,
	`
	create index comment_post_id_idx on comments (post_id);
	`,
	`
	create index comment_post_id_created_idx on comments (post_id, created);
	`,
	`
	create index comment_created_idx on comments (created);
	`,
}

const SQLForCountPosts = `SELECT COUNT(*) FROM posts`
const SQLForCountComments = `SELECT COUNT(*) FROM comments`

const SQLForInsertPost = `INSERT INTO posts (title, content, stats) VALUES (?1, ?2, ?3)`
const SQLForInsertPostWithCreated = `INSERT INTO posts (title, content, stats, created) VALUES (?1, ?2, ?3, ?4)`
const SQLForInsertComment = `INSERT INTO comments (post_id, name, content, stats) VALUES (?1, ?2, ?3, ?4)`
const SQLForInsertCommentWithCreated = `INSERT INTO comments (post_id, name, content, stats, created) VALUES (?1, ?2, ?3, ?4, ?5)`

const SQLForSelectPostByID = `SELECT title, content, created, stats FROM posts WHERE id = ?1`
const SQLForSelectCommentsByPostID = `SELECT id, name, content, created, stats FROM comments WHERE post_id = ?1 ORDER BY created`

const SQLForQueryCorrelated = `
	SELECT
		id,
		title,
		(SELECT COUNT(*) FROM comments WHERE post_id = posts.id) as comment_count,
		(SELECT AVG(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS avg_comment_length,
		(SELECT MAX(LENGTH(content)) FROM comments WHERE post_id = posts.id) AS max_comment_length
	FROM posts
	`
const SQLForQueryGroupBy = `
	SELECT
		strftime('%Y-%m', created) AS month,
		COUNT(*) as month_total
	FROM posts
	GROUP BY month
	ORDER BY month
	`
const SQLForQueryJSON = `
	SELECT
		date(created) as day,
		SUM(json_extract(stats, '$.lorem')) as sum_lorem,
		AVG(json_extract(stats, '$.ipsum.dolor')) as avg_dolor,
		MAX(json_extract(stats, '$.lorem.sit')) as max_sit
	FROM posts
	GROUP BY day
	ORDER BY day
	`
const SQLForQueryNonRecursiveCTE = `
	WITH day_totals AS (
		SELECT date(created) as day, COUNT(*) as day_total
		FROM posts
		GROUP BY day
	)
	SELECT day, day_total,
		SUM(day_total) OVER (ORDER BY day) as running_total
	FROM day_totals
	ORDER BY day
	`
const SQLForQueryOrderBy = `
	SELECT
		name, created, id
	FROM comments
	ORDER BY name, created, id
	`
const SQLForQueryRecursiveCTE = `
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
	`
