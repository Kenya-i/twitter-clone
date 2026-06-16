package repository

import (
	"context"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type followRepository struct {
	db *pgxpool.Pool
}

func NewFollowRepository(db *pgxpool.Pool) domain.FollowRepository {
	return &followRepository{db: db}
}

func (r *followRepository) Create(follow *domain.Follow) error {
	query := `
		INSERT INTO follows (follower_id, following_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING id`

	follow.CreatedAt = time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.QueryRow(ctx, query, follow.FollowerID, follow.FollowingID, follow.CreatedAt).Scan(&follow.ID)
}

func (r *followRepository) Delete(followerID, followingID string) error {
	query := `DELETE FROM follows WHERE follower_id = $1 AND following_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx, query, followerID, followingID)
	return err
}

func (r *followRepository) Exists(followerID, followingID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM follows WHERE follower_id = $1 AND following_id = $2)`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var exists bool
	err := r.db.QueryRow(ctx, query, followerID, followingID).Scan(&exists)
	return exists, err
}

func (r *followRepository) CountFollowers(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE following_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *followRepository) CountFollowing(userID string) (int, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var count int
	err := r.db.QueryRow(ctx, query, userID).Scan(&count)
	return count, err
}

func (r *followRepository) GetFollowingIDs(userID string) ([]string, error) {
	query := `SELECT following_id FROM follows WHERE follower_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	return ids, rows.Err()
}
