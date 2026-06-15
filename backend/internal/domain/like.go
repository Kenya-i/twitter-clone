package domain

import "time"

type Like struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TweetID   string    `json:"tweet_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LikeRepository interface {
	Create(like *Like) error
	Delete(userID, tweetID string) error
	Exists(userID, tweetID string) (bool, error)
	CountByTweetID(tweetID string) (int, error)
}
