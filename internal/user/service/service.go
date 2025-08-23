package service

import (
	responseDto "bulletin-board/internal/ad/dto"
	"bulletin-board/internal/user"
	"bulletin-board/internal/user/dto"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const signingKey = "iuNvi8sa5oiHOajKfn93hFL93gb"

type Service struct {
	repository user.Repository
}

type TokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func NewService(repository user.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context) ([]dto.ResponseUser, error) {
	users, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	responseUsers := make([]dto.ResponseUser, 0)
	for _, usr := range users {
		responseUser := dto.ToDto(usr)
		responseUsers = append(responseUsers, responseUser)
	}
	return responseUsers, nil
}

func (s *Service) GetByID(ctx context.Context, id int) (dto.ResponseUser, error) {
	if id < 1 {
		return dto.ResponseUser{}, user.ErrInvalidUserId
	}
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return dto.ResponseUser{}, err
	}
	return dto.ToDto(user), nil
}

func (s *Service) GetUsersAds(ctx context.Context, userId int) ([]responseDto.ResponseAd, error) {
	if userId < 1 {
		return []responseDto.ResponseAd{}, user.ErrInvalidUserId
	}
	ads, err := s.repository.GetUsersAds(ctx, userId)
	if err != nil {
		return nil, err
	}
	responseAds := make([]responseDto.ResponseAd, 0)

	for _, ad := range ads {
		responseAd := responseDto.ToDto(ad)
		responseAds = append(responseAds, responseAd)
	}
	return responseAds, nil
}

func (s *Service) Create(ctx context.Context, newUser dto.RequestUser) (dto.ResponseUser, error) {
	user := dto.ToUser(newUser)
	hash, err := s.generatePasswordHash(user.Password)
	if err != nil {
		return dto.ResponseUser{}, err
	}
	user.Password = hash
	createdUser, err := s.repository.Create(ctx, user)
	if err != nil {
		return dto.ResponseUser{}, err
	}
	return dto.ToDto(createdUser), nil
}

func (s *Service) Update(ctx context.Context, requestUser dto.RequestUser, id int) (dto.ResponseUser, error) {
	if id < 1 {
		return dto.ResponseUser{}, user.ErrInvalidUserId
	}
	user := dto.ToUser(requestUser)
	updatedUser, err := s.repository.Update(ctx, user, id)
	if err != nil {
		return dto.ResponseUser{}, err
	}
	return dto.ToDto(updatedUser), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 1 {
		return user.ErrInvalidUserId
	}
	return s.repository.Delete(ctx, id)
}

func (s *Service) generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (s *Service) GenerateToken(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("user not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		}, user.ID,
	})

	return token.SignedString([]byte(signingKey))
}
