package domain

import "time"

type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TweetRepository interface {
	Create(tweet *Tweet) error
	FindByID(id string) (*Tweet, error)
	FindAll() ([]*Tweet, error)
	Delete(id string) error
}

type TweetUsecase interface {
	Post(userID, content string) (*Tweet, error)
	GetTweet(id string) (*Tweet, error)
	GetTimeline() ([]*Tweet, error)
	Delete(userID, tweetID string) error
}
