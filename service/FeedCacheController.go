package service

import (
	"HighArch/storage"
	"errors"
	"github.com/go-redis/redis"
	"time"
)

type FeedCacheController interface {
	InvalidateFeedCacheForUser(userId string)
	InvalidateFeedsCacheForFriends(userId string)
	ListenHandleFeedUpdate()
}

type redisFeedCacheController struct {
	redisClient      *redis.Client
	postsStore       storage.PostsStore
	postsCacheStore  storage.PostsCacheStore
	friendLinksStore storage.FriendLinksStore
}

func NewRedisCacheController(redisClient *redis.Client, postStore storage.PostsStore, postCacheStore storage.PostsCacheStore, friendLinksStore storage.FriendLinksStore) *redisFeedCacheController {
	return &redisFeedCacheController{
		redisClient:      redisClient,
		postsStore:       postStore,
		postsCacheStore:  postCacheStore,
		friendLinksStore: friendLinksStore,
	}
}

func (c *redisFeedCacheController) InvalidateFeedCacheForUser(userId string) {
	_, err := c.redisClient.RPush(topFeedQueue, userId).Result()
	if err != nil {
		print(err)
	}
	println("invalidateFeedCacheForUser: " + userId)
}

func (c *redisFeedCacheController) InvalidateFeedsCacheForFriends(userId string) {
	friendsIds, err := c.friendLinksStore.GetFriendsIds(userId)
	if err != nil {
		println(err)
		return
	}
	for _, friendId := range friendsIds {
		c.InvalidateFeedCacheForUser(friendId)
	}
}

func (c *redisFeedCacheController) ListenHandleFeedUpdate() {
	for {
		userId, err := c.redisClient.RPop(topFeedQueue).Result()
		if errors.Is(err, redis.Nil) {
			// Wait before popping the next item
			time.Sleep(time.Duration(50) * time.Millisecond)
			continue
		}
		if err != nil {
			println("ListenHandleFeedUpdate error:", err)
			continue
		}
		if c.postsCacheStore.HasTopFeed(userId) {
			// Update the top feed only for existed cache
			posts, err := c.postsStore.GetFeed(userId, 0, 30)
			if err != nil {
				print(err)
				c.postsCacheStore.RemoveTopFeed(userId)
			} else {
				println("ListenHandleFeedUpdate SetTopFeed: " + userId)
				err = c.postsCacheStore.SetTopFeed(userId, posts)
				if err != nil {
					print(err)
				}
			}
		}
	}
}

const topFeedQueue = "top_feed_queue"
