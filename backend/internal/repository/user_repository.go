package repository

import (
	"context"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, hashed_password, display_name, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.HashedPassword,
		user.DisplayName,
		user.Bio,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	query := `SELECT id, username, email, hashed_password, display_name, bio, created_at, updated_at FROM users WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.DisplayName,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	query := `SELECT id, username, email, hashed_password, display_name, bio, created_at, updated_at FROM users WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.DisplayName,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, email, hashed_password, display_name, bio, created_at, updated_at FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user domain.User
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.HashedPassword,
		&user.DisplayName,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
