package repository

import (
	"context"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type likeRepository struct {
	db *pgxpool.Pool
}

func NewLikeRepository(db *pgxpool.Pool) domain.LikeRepository {
	return &likeRepository{db: db}
}

func (r *likeRepository) Create(like *domain.Like) error {
	query := `
		INSERT INTO likes (user_id, tweet_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	like.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.QueryRow(ctx, query, like.UserID, like.TweetID, like.CreatedAt).Scan(&like.ID)
}

func (r *likeRepository) Delete(userID, tweetID string) error {
	query := `DELETE FROM likes WHERE user_id = $1 AND tweet_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, userID, tweetID)
	return err
}

func (r *likeRepository) Exists(userID, tweetID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM likes WHERE user_id = $1 AND tweet_id = $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, tweetID).Scan(&exists)
	return exists, err
}

func (r *likeRepository) CountByTweetID(tweetID string) (int, error) {
	query := `SELECT COUNT(*) FROM likes WHERE tweet_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRow(ctx, query, tweetID).Scan(&count)
	return count, err
}
