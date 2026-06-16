package domain

import "time"

type Follow struct {
	ID          string    `json:"id"`
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type FollowRepository interface {
	Create(follow *Follow) error
	Delete(followerID, followingID string) error
	Exists(followerID, followingID string) (bool, error)
	CountFollowers(userID string) (int, error)
	CountFollowing(userID string) (int, error)
	GetFollowingIDs(userID string) ([]string, error)
}

type FollowUsecase interface {
	Follow(followerID, followingID string) error
	Unfollow(followerID, followingID string) error
	IsFollowing(followerID, followingID string) (bool, error)
	GetFollowCounts(userID string) (followers int, following int, err error)
}
