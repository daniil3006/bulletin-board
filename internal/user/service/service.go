package service

import (
	responseDto "bulletin-board/internal/ad/dto"
	"bulletin-board/internal/user"
	"bulletin-board/internal/user/dto"
	"context"
)

type Service struct {
	repository user.Repository
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
	_, err := s.repository.Create(ctx, user)
	if err != nil {
		return dto.ResponseUser{}, nil
	}
	return dto.ToDto(user), nil
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
