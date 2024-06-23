package service

import (
	"HighArch/api"
	"HighArch/entity"
	"HighArch/storage"
	"time"
)

type FeedService struct {
	postStore        storage.PostsStore
	postsCacheStore  storage.PostsCacheStore
	friendLinksStore storage.FriendLinksStore
}

func NewFeedService(postStore storage.PostsStore, postsCacheStore storage.PostsCacheStore, friendLinksStore storage.FriendLinksStore) *FeedService {
	return &FeedService{
		postStore:        postStore,
		postsCacheStore:  postsCacheStore,
		friendLinksStore: friendLinksStore,
	}
}

func (s *FeedService) GetFeed(userId string, offset, limit int) ([]api.PostApiModel, error) {
	var posts []entity.Post = nil
	var err error = nil

	// check top feed in cache
	if isCacheTopFeed(offset, limit) {
		cachedPosts, err := s.postsCacheStore.GetTopFeed(userId)
		if err == nil && cachedPosts != nil {
			println("Top Feed Cache Hit")
			posts = cachedPosts
		}
		if err != nil {
			println(err)
		}
	}

	if posts == nil { // no cache
		posts, err = s.postStore.GetFeed(userId, offset, limit)
		if err != nil {
			return nil, ErrorStoreError
		}
	}
	result := make([]api.PostApiModel, 0)
	for _, post := range posts {
		// TODO move to mapper
		result = append(result, api.PostApiModel{
			Id:         post.Id,
			Text:       post.Text,
			AuthorId:   post.AuthorId,
			CreateTime: formatUnixTimestampToString(post.CreateTime, time.DateTime),
		})
	}

	if isCacheTopFeed(offset, limit) {
		println("Put cache")
		err = s.postsCacheStore.SetTopFeed(userId, posts)
		if err != nil {
			print(err)
		}
	}
	return result, nil
}

// TODO the same constants (0 and 30) are used in FeedCacheController code
func isCacheTopFeed(offset, limit int) bool {
	return offset == 0 && limit == 30
}
