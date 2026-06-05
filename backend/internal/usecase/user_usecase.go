package usecase

import (
	"errors"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(userRepo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Register(username, email, password, displayName string) (*domain.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:       username,
		Email:          email,
		HashedPassword: string(hashed),
		DisplayName:    displayName,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)); err != nil {
		return "", errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	return user.ID, nil
}

func (u *userUsecase) GetProfile(id string) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}
