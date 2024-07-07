package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"time"
)

type PostService struct {
	postStore           storage.PostsStore
	feedCacheController FeedCacheController
	feedWsController    FeedWsController
}

func NewPostService(postStore storage.PostsStore, feedCacheController FeedCacheController, feedWsController FeedWsController) *PostService {
	return &PostService{
		postStore:           postStore,
		feedCacheController: feedCacheController,
		feedWsController:    feedWsController,
	}
}

func (s *PostService) CreatePost(postText string, authorId string) (*api.PostCreateSuccessApiModel, error) {
	var err = validatePost(postText)
	if err != nil {
		return nil, err
	}
	newPost := entity.Post{
		Text:       postText,
		AuthorId:   authorId,
		CreateTime: time.Now().UnixMilli(),
	}
	id, err := s.postStore.CreatePost(newPost)
	if err != nil {
		return nil, ErrorStoreError
	}

	go s.feedCacheController.InvalidateFeedsCacheForFriends(authorId)
	go s.feedWsController.HandleNewPostCreated(authorId, newPost)

	return &api.PostCreateSuccessApiModel{PostId: *id}, nil
}

func (s *PostService) GetPost(id string) (*api.PostApiModel, error) {
	var post, err = s.postStore.GetPost(id)
	if err != nil {
		return nil, ErrorStoreError
	}
	if post == nil {
		return nil, ErrorNotFound
	}

	// TODO move to mapper
	var result = api.PostApiModel{
		Id:         post.Id,
		Text:       post.Text,
		AuthorId:   post.AuthorId,
		CreateTime: formatUnixTimestampToString(post.CreateTime, time.DateTime),
	}

	return &result, nil
}

func validatePost(postText string) error {
	if len(postText) <= 0 {
		return ErrorValidation
	}

	return nil
}
