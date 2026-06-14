package repository

import (
	"context"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tweetRepository struct {
	db *pgxpool.Pool
}

func NewTweetRepository(db *pgxpool.Pool) domain.TweetRepository {
	return &tweetRepository{db: db}
}

func (r *tweetRepository) Create(tweet *domain.Tweet) error {
	query := `
		INSERT INTO tweets (user_id, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	now := time.Now()
	tweet.CreatedAt = now
	tweet.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.QueryRow(ctx, query,
		tweet.UserID,
		tweet.Content,
		tweet.CreatedAt,
		tweet.UpdatedAt,
	).Scan(&tweet.ID)
}

func (r *tweetRepository) FindByID(id string) (*domain.Tweet, error) {
	query := `SELECT id, user_id, content, created_at, updated_at FROM tweets WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var tweet domain.Tweet
	err := r.db.QueryRow(ctx, query, id).Scan(
		&tweet.ID,
		&tweet.UserID,
		&tweet.Content,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &tweet, nil
}

func (r *tweetRepository) FindAll() ([]*domain.Tweet, error) {
	query := `SELECT id, user_id, content, created_at, updated_at FROM tweets ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tweets := []*domain.Tweet{}
	for rows.Next() {
		var tweet domain.Tweet
		if err := rows.Scan(
			&tweet.ID,
			&tweet.UserID,
			&tweet.Content,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tweets = append(tweets, &tweet)
	}

	return tweets, rows.Err()
}

func (r *tweetRepository) Delete(id string) error {
	query := `DELETE FROM tweets WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, id)
	return err
}
