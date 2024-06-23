package storage

import (
	"HighArch/entity"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis"
)

type PostsCacheStore interface {
	HasTopFeed(userId string) bool
	GetTopFeed(userId string) ([]entity.Post, error)
	SetTopFeed(userId string, posts []entity.Post) error
	RemoveTopFeed(userId string) error
}

type RedisPostsCacheStore struct {
	redisClient *redis.Client
}

func NewRedisPostsCacheStore(client *redis.Client) *RedisPostsCacheStore {
	return &RedisPostsCacheStore{redisClient: client}
}

func (s *RedisPostsCacheStore) HasTopFeed(userId string) bool {
	exists, _ := s.redisClient.Exists(getTopFeedKey(userId)).Result()
	return exists == 1
}

func (s *RedisPostsCacheStore) GetTopFeed(userId string) ([]entity.Post, error) {
	serializedPosts, err := s.redisClient.Get(getTopFeedKey(userId)).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var posts []entity.Post
	err = json.Unmarshal([]byte(serializedPosts), &posts)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *RedisPostsCacheStore) SetTopFeed(userId string, posts []entity.Post) error {
	serializedPosts, err := json.Marshal(posts)
	if err != nil {
		return err
	}
	_, err = s.redisClient.Set(getTopFeedKey(userId), serializedPosts, 0).Result()
	if err != nil {
		s.redisClient.Del(getTopFeedKey(userId))
		return err
	}
	return nil
}

func (s *RedisPostsCacheStore) RemoveTopFeed(userId string) error {
	return s.redisClient.Del(getTopFeedKey(userId)).Err()
}

func getTopFeedKey(userId string) string {
	return "TopFeed:" + userId
}
