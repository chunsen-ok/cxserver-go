package writer

type PostTag struct {
	PostID int `json:"post_id"`
	TagID  int `json:"tag_id"`
}

const PostTagSQL = `
CREATE TABLE IF NOT EXISTS writer.post_tags (
	post_id integer NOT NULL,
	tag_id integer NOT NULL,
	UNIQUE (post_id, tag_id)
);`
