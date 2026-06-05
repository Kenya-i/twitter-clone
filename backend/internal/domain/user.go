package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username       string             `bson:"username" json:"username"`
	Email          string             `bson:"email" json:"email"`
	HashedPassword string             `bson:"hashed_password" json:"-"`
	DisplayName    string             `bson:"display_name" json:"display_name"`
	Bio            string             `bson:"bio" json:"bio"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
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
