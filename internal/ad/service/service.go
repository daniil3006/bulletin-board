package service

import (
	"bulletin-board/internal/ad"
	"bulletin-board/internal/ad/dto"
	"context"
	"errors"
)

type Service struct {
	repository ad.Repository
}

func NewService(repository ad.Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) GetAll(ctx context.Context) ([]dto.ResponseAd, error) {
	ads, err := s.repository.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	responseAds := make([]dto.ResponseAd, 0)
	for _, ad := range ads {
		responseAd := dto.ToDto(ad)
		responseAds = append(responseAds, responseAd)
	}
	return responseAds, nil
}

func (s *Service) GetByID(ctx context.Context, ID int) (dto.ResponseAd, error) {
	if ID <= 0 {
		return dto.ResponseAd{}, errors.New("invalid id")
	}
	ad, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return dto.ResponseAd{}, err
	}
	return dto.ToDto(ad), nil
}

func (s *Service) Create(ctx context.Context, requestAd dto.RequestAd) (dto.ResponseAd, error) {
	ad := dto.ToAd(requestAd)
	if err := checkValidityAd(ad); err != nil {
		return dto.ResponseAd{}, err
	}
	ad, err := s.repository.Create(ctx, ad)
	if err != nil {
		return dto.ResponseAd{}, err
	}
	return dto.ToDto(ad), nil
}

func (s *Service) Update(ctx context.Context, requestAd dto.RequestAd, id int) (dto.ResponseAd, error) {
	reqAd := dto.ToAd(requestAd)

	authId, ok := ctx.Value("user_id").(int)
	if !ok {
		return dto.ResponseAd{}, errors.New("invalid auth")
	}

	if !s.checkValidityUser(ctx, authId, id) {
		return dto.ResponseAd{}, ad.ErrForbidden
	}

	if err := checkValidityAd(reqAd); err != nil {
		return dto.ResponseAd{}, err
	}

	reqAd, err := s.repository.Update(ctx, reqAd, id)
	if err != nil {
		return dto.ResponseAd{}, err
	}

	reqAd.UserID = authId

	return dto.ToDto(reqAd), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id < 0 {
		return errors.New("invalid id")
	}
	return s.repository.Delete(ctx, id)
}

func checkValidityAd(ad ad.Ad) error {
	if ad.Price < 0 {
		return errors.New("invalid price")
	}
	if ad.Title == "" {
		return errors.New("invalid title")
	}

	return nil
}

func (s *Service) checkValidityUser(ctx context.Context, authUser, adId int) bool {
	ad, _ := s.repository.GetByID(ctx, adId)
	return ad.UserID == authUser
}
