package go_sqlite_bench

import _ "embed"

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

//go:embed gsb_common_sql_query_correlated.sql
var SQLForQueryCorrelated string

//go:embed gsb_common_sql_query_groupby.sql
var SQLForQueryGroupBy string

//go:embed gsb_common_sql_query_json.sql
var SQLForQueryJSON string

//go:embed gsb_common_sql_query_nonrecursivecte.sql
var SQLForQueryNonRecursiveCTE string

//go:embed gsb_common_sql_query_orderby.sql
var SQLForQueryOrderBy string

//go:embed gsb_common_sql_query_recursivecte.sql
var SQLForQueryRecursiveCTE string
