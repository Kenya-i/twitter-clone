package usecase

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/Kenya-i/twitter-clone/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo     domain.UserRepository
	imageStorage domain.ImageStorage
	jwtSecret    string
}

func NewUserUsecase(userRepo domain.UserRepository, imageStorage domain.ImageStorage, jwtSecret string) domain.UserUsecase {
	return &userUsecase{userRepo: userRepo, imageStorage: imageStorage, jwtSecret: jwtSecret}
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString([]byte(u.jwtSecret))
}

func (u *userUsecase) GetProfile(id string) (*domain.User, error) {
	return u.userRepo.FindByID(id)
}

func (u *userUsecase) GetUsers() ([]*domain.User, error) {
	return u.userRepo.FindAll()
}

func (u *userUsecase) UpdateAvatar(userID string, file io.Reader, filename string, contentType string) (*domain.User, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(filename)
	key := fmt.Sprintf("avatars/%s-%d%s", userID, time.Now().UnixNano(), ext)

	url, err := u.imageStorage.Upload(file, key, contentType)
	if err != nil {
		return nil, err
	}

	user.AvatarURL = url

	if err := u.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
