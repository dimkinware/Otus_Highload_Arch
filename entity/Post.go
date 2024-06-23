package entity

type Post struct {
	Id         string `db:"id"`
	Text       string `db:"post_text"`
	AuthorId   string `db:"author_user_id"`
	CreateTime int64  `db:"create_time"`
}
