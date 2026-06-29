package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/redis/go-redis/v9"
)

type cachedTweetRepository struct {
	inner domain.TweetRepository
	redis *redis.Client
}

func NewCachedTweetRepository(inner domain.TweetRepository, redisClient *redis.Client) domain.TweetRepository {
	return &cachedTweetRepository{inner: inner, redis: redisClient}
}

func (r *cachedTweetRepository) Create(tweet *domain.Tweet) error {
	if err := r.inner.Create(tweet); err != nil {
		return err
	}

	ctx := context.Background()
	key := fmt.Sprintf("timeline:%s:20", tweet.UserID)
	r.redis.Del(ctx, key)

	return nil
}

func (r *cachedTweetRepository) FindByID(id string) (*domain.Tweet, error) {
	return r.inner.FindByID(id)
}

func (r *cachedTweetRepository) FindByFollowing(userID string, cursor *time.Time, limit int) ([]*domain.Tweet, error) {
	// カーソルがある（2ページ目以降）はキャッシュ対象外。最新ページだけキャッシュする
	if cursor != nil {
		return r.inner.FindByFollowing(userID, cursor, limit)
	}

	ctx := context.Background()
	key := fmt.Sprintf("timeline:%s:%d", userID, limit)

	if cached, err := r.redis.Get(ctx, key).Result(); err == nil {
		var tweets []*domain.Tweet
		if jsonErr := json.Unmarshal([]byte(cached), &tweets); jsonErr == nil {
			return tweets, nil
		}
	}

	tweets, err := r.inner.FindByFollowing(userID, cursor, limit)
	if err != nil {
		return nil, err
	}

	if data, jsonErr := json.Marshal(tweets); jsonErr == nil {
		r.redis.Set(ctx, key, data, 30*time.Second)
	}

	return tweets, nil
}

func (r *cachedTweetRepository) Search(query string, cursor *time.Time, limit int) ([]*domain.Tweet, error) {
	return r.inner.Search(query, cursor, limit)
}

func (r *cachedTweetRepository) Update(tweet *domain.Tweet) error {
	return r.inner.Update(tweet)
}

func (r *cachedTweetRepository) Delete(id string) error {
	return r.inner.Delete(id)
}
