package domain

import "time"

type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	LikeCount int       `json:"like_count"`
	LikedByMe bool      `json:"liked_by_me"`
}

type TweetRepository interface {
	Create(tweet *Tweet) error
	FindByID(id string) (*Tweet, error)
	FindAll() ([]*Tweet, error)
	Update(tweet *Tweet) error
	Delete(id string) error
}

type TweetUsecase interface {
	Post(userID, content string) (*Tweet, error)
	GetTweet(id, userID string) (*Tweet, error)
	GetTimeline(userID string) ([]*Tweet, error)
	Update(userID, tweetID, content string) (*Tweet, error)
	Delete(userID, tweetID string) error
	Like(userID, tweetID string) error
	Unlike(userID, tweetID string) error
}
