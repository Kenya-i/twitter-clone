package domain

import "time"

type User struct {
	ID             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	DisplayName    string    `json:"display_name"`
	Bio            string    `json:"bio"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	FindByID(id string) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByUsername(username string) (*User, error)
}

type UserUsecase interface {
	Register(username, email, password, displayName string) (*User, error)
	Login(email, password string) (string, error)
	GetProfile(id string) (*User, error)
}
