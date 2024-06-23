package api

type PostApiModel struct {
	Id         string `json:"id"`
	Text       string `json:"text"`
	AuthorId   string `json:"author_user_id"`
	CreateTime string `json:"create_time"`
}

type PostCreateApiModel struct {
	Text string `json:"text"`
}

type PostCreateSuccessApiModel struct {
	PostId string `json:"post_id"`
}
