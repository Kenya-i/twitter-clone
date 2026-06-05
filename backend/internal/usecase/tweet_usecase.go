package usecase

import (
	"errors"

	"github.com/Kenya-i/twitter-clone/internal/domain"
)

type tweetUsecase struct {
	tweetRepo domain.TweetRepository
}

func NewTweetUsecase(tweetRepo domain.TweetRepository) domain.TweetUsecase {
	return &tweetUsecase{tweetRepo: tweetRepo}
}

func (u *tweetUsecase) Post(userID, content string) (*domain.Tweet, error) {
	tweet := &domain.Tweet{
		UserID:  userID,
		Content: content,
	}

	if err := u.tweetRepo.Create(tweet); err != nil {
		return nil, err
	}

	return tweet, nil
}

func (u *tweetUsecase) GetTweet(id string) (*domain.Tweet, error) {
	return u.tweetRepo.FindByID(id)
}

func (u *tweetUsecase) Delete(userID, tweetID string) error {
	tweet, err := u.tweetRepo.FindByID(tweetID)
	if err != nil {
		return errors.New("ツイートが見つかりません")
	}

	if tweet.UserID != userID {
		return errors.New("削除権限がありません")
	}

	return u.tweetRepo.Delete(tweetID)
}
