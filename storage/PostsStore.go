package storage

import (
	"HighArch/entity"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PostsStore interface {
	CreatePost(post entity.Post) (*string, error)
	GetPost(id string) (*entity.Post, error)
	GetFeed(userId string, offset, limit int) ([]entity.Post, error)

	// TODO DeletePost & UpdatePost
}

type dbPostsStore struct {
	db *sqlx.DB
}

func NewDbPostsStore(db *sqlx.DB) PostsStore {
	return &dbPostsStore{
		db: db,
	}
}

func (s dbPostsStore) CreatePost(post entity.Post) (*string, error) {
	var postId = post.Id
	if len(postId) <= 0 {
		postId = uuid.NewString()
	}
	query := "INSERT INTO posts(id, author_user_id, post_text, create_time) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, postId, post.AuthorId, post.Text, post.CreateTime)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &postId, nil
}

func (s dbPostsStore) GetPost(id string) (*entity.Post, error) {
	rows, err := s.db.Queryx("SELECT * FROM posts WHERE id = $1 limit 1", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var post entity.Post
		err = rows.StructScan(&post)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &post, nil
	}

	return nil, nil
}

func (s dbPostsStore) GetFeed(userId string, offset, limit int) ([]entity.Post, error) {
	query := "SELECT posts.* FROM posts where (author_user_id in (select user_id_f1 from friends where user_id_f2 = $1) OR author_user_id in (select user_id_f2 from friends where user_id_f1 = $1)) AND author_user_id != $1 ORDER BY create_time desc LIMIT $2 OFFSET $3;"
	rows, err := s.db.Queryx(query, userId, limit, offset)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	var posts []entity.Post
	for rows.Next() {
		var post entity.Post
		err = rows.StructScan(&post)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
