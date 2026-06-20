package usecase

import (
	"errors"
	"testing"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
)

// mockTweetRepository は domain.TweetRepository を実装する偽のリポジトリ
type mockTweetRepository struct {
	tweets map[string]*domain.Tweet
}

func (m *mockTweetRepository) Create(tweet *domain.Tweet) error {
	return nil
}

func (m *mockTweetRepository) FindByID(id string) (*domain.Tweet, error) {
	tweet, ok := m.tweets[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return tweet, nil
}

func (m *mockTweetRepository) FindByFollowing(userID string, cursor *time.Time, limit int) ([]*domain.Tweet, error) {
	return nil, nil
}

func (m *mockTweetRepository) Update(tweet *domain.Tweet) error {
	return nil
}

func (m *mockTweetRepository) Delete(id string) error {
	delete(m.tweets, id)
	return nil
}

// mockLikeRepository は domain.LikeRepository を実装する偽のリポジトリ
// Delete のテストでは使わないが、tweetUsecase の生成に必要
type mockLikeRepository struct{}

func (m *mockLikeRepository) Create(like *domain.Like) error              { return nil }
func (m *mockLikeRepository) Delete(userID, tweetID string) error         { return nil }
func (m *mockLikeRepository) Exists(userID, tweetID string) (bool, error) { return false, nil }
func (m *mockLikeRepository) CountByTweetID(tweetID string) (int, error)  { return 0, nil }

func TestDelete_他人のツイートは削除できない(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		tweets: map[string]*domain.Tweet{
			"tweet1": {ID: "tweet1", UserID: "owner-user"},
		},
	}
	likeRepo := &mockLikeRepository{}

	u := NewTweetUsecase(tweetRepo, likeRepo)

	err := u.Delete("other-user", "tweet1")

	if err == nil {
		t.Fatal("他人のツイートなのにエラーが返らなかった")
	}
}

func TestDelete_自分のツイートは削除できる(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		tweets: map[string]*domain.Tweet{
			"tweet1": {ID: "tweet1", UserID: "owner-user"},
		},
	}
	likeRepo := &mockLikeRepository{}

	u := NewTweetUsecase(tweetRepo, likeRepo)

	err := u.Delete("owner-user", "tweet1")

	if err != nil {
		t.Fatalf("自分のツイートなのにエラーが返った: %v", err)
	}
}

func TestUpdate_他人のツイートは編集できない(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		tweets: map[string]*domain.Tweet{
			"tweet1": {ID: "tweet1", UserID: "owner-user", Content: "元の内容"},
		},
	}
	likeRepo := &mockLikeRepository{}

	u := NewTweetUsecase(tweetRepo, likeRepo)

	_, err := u.Update("other-user", "tweet1", "書き換えたい内容")

	if err == nil {
		t.Fatal("他人のツイートなのにエラーが返らなかった")
	}
}

func TestUpdate_自分のツイートは編集できる(t *testing.T) {
	tweetRepo := &mockTweetRepository{
		tweets: map[string]*domain.Tweet{
			"tweet1": {ID: "tweet1", UserID: "owner-user", Content: "元の内容"},
		},
	}
	likeRepo := &mockLikeRepository{}

	u := NewTweetUsecase(tweetRepo, likeRepo)

	updated, err := u.Update("owner-user", "tweet1", "新しい内容")

	if err != nil {
		t.Fatalf("自分のツイートなのにエラーが返った: %v", err)
	}

	if updated.Content != "新しい内容" {
		t.Errorf("内容が更新されていない: got %q", updated.Content)
	}
}
