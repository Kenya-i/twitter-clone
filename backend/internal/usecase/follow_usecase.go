package usecase

import (
	"errors"

	"github.com/Kenya-i/twitter-clone/internal/domain"
)

type followUsecase struct {
	followRepo domain.FollowRepository
}

func NewFollowUsecase(followRepo domain.FollowRepository) domain.FollowUsecase {
	return &followUsecase{followRepo: followRepo}
}

func (u *followUsecase) Follow(followerID, followingID string) error {
	if followerID == followingID {
		return errors.New("自分自身をフォローすることはできません")
	}

	exists, err := u.followRepo.Exists(followerID, followingID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	follow := &domain.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
	}

	return u.followRepo.Create(follow)
}

func (u *followUsecase) Unfollow(followerID, followingID string) error {
	return u.followRepo.Delete(followerID, followingID)
}

func (u *followUsecase) IsFollowing(followerID, followingID string) (bool, error) {
	return u.followRepo.Exists(followerID, followingID)
}

func (u *followUsecase) GetFollowCounts(userID string) (int, int, error) {
	followers, err := u.followRepo.CountFollowers(userID)
	if err != nil {
		return 0, 0, err
	}

	following, err := u.followRepo.CountFollowing(userID)
	if err != nil {
		return 0, 0, err
	}

	return followers, following, nil
}
