package usecase

import (
	"errors"

	"github.com/Kenya-i/twitter-clone/internal/domain"
)

type tweetUsecase struct {
	tweetRepo domain.TweetRepository
	likeRepo  domain.LikeRepository
}

func NewTweetUsecase(tweetRepo domain.TweetRepository, likeRepo domain.LikeRepository) domain.TweetUsecase {
	return &tweetUsecase{tweetRepo: tweetRepo, likeRepo: likeRepo}
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

func (u *tweetUsecase) GetTweet(id, userID string) (*domain.Tweet, error) {
	tweet, err := u.tweetRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if err := u.attachLikeInfo(tweet, userID); err != nil {
		return nil, err
	}

	return tweet, nil
}

func (u *tweetUsecase) GetTimeline(userID string) ([]*domain.Tweet, error) {
	tweets, err := u.tweetRepo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, tweet := range tweets {
		if err := u.attachLikeInfo(tweet, userID); err != nil {
			return nil, err
		}
	}

	return tweets, nil
}

func (u *tweetUsecase) attachLikeInfo(tweet *domain.Tweet, userID string) error {
	count, err := u.likeRepo.CountByTweetID(tweet.ID)
	if err != nil {
		return err
	}
	tweet.LikeCount = count

	liked, err := u.likeRepo.Exists(userID, tweet.ID)
	if err != nil {
		return err
	}
	tweet.LikedByMe = liked

	return nil
}

func (u *tweetUsecase) Update(userID, tweetID, content string) (*domain.Tweet, error) {
	tweet, err := u.tweetRepo.FindByID(tweetID)
	if err != nil {
		return nil, errors.New("ツイートが見つかりません")
	}

	if tweet.UserID != userID {
		return nil, errors.New("編集権限がありません")
	}

	tweet.Content = content

	if err := u.tweetRepo.Update(tweet); err != nil {
		return nil, err
	}

	return tweet, nil
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

func (u *tweetUsecase) Like(userID, tweetID string) error {
	exists, err := u.likeRepo.Exists(userID, tweetID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	like := &domain.Like{
		UserID:  userID,
		TweetID: tweetID,
	}

	return u.likeRepo.Create(like)
}

func (u *tweetUsecase) Unlike(userID, tweetID string) error {
	return u.likeRepo.Delete(userID, tweetID)
}
