package service

import (
	"bulletin-board/internal/ad"
	"bulletin-board/internal/ad/dto"
	"bulletin-board/internal/redisdb"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type Service struct {
	repository ad.Repository
	rds        redisdb.RedisClient
}

func NewService(repository ad.Repository, rds redisdb.RedisClient) *Service {
	return &Service{repository: repository, rds: rds}
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

	adObj, err := s.getAdFromRedis(ctx, ID)
	if err != nil {
		log.Printf("Redis error: %v", err)
	}

	if adObj.ID != 0 {
		return dto.ToDto(adObj), nil
	}

	ad, err := s.repository.GetByID(ctx, ID)
	if err != nil {
		return dto.ResponseAd{}, err
	}

	err = s.addToRedis(ctx, ID, ad, 10*time.Minute)
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

	reqAd.UserID = authId

	reqAd, err := s.repository.Update(ctx, reqAd, id)
	if err != nil {
		return dto.ResponseAd{}, err
	}

	s.deleteFromRedis(ctx, id)
	return dto.ToDto(reqAd), nil
}

func (s *Service) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("invalid id")
	}

	authId, ok := ctx.Value("user_id").(int)

	if !ok {
		return errors.New("invalid auth")
	}

	if !s.checkValidityUser(ctx, authId, id) {
		return ad.ErrForbidden
	}

	s.deleteFromRedis(ctx, id)
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
	ad, err := s.repository.GetByID(ctx, adId)
	if err != nil {
		return false
	}
	return ad.UserID == authUser
}

func (s *Service) addToRedis(ctx context.Context, ID int, adObj ad.Ad, tm time.Duration) error {
	jsonAd, err := json.Marshal(adObj)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("ad:%d", ID)
	return s.rds.Rds.Set(ctx, key, jsonAd, tm).Err()
}

func (s *Service) getAdFromRedis(ctx context.Context, ID int) (ad.Ad, error) {
	key := fmt.Sprintf("ad:%d", ID)
	jsonAd, err := s.rds.Rds.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return ad.Ad{}, nil
		}
		return ad.Ad{}, err
	}

	var adObj ad.Ad
	err = json.Unmarshal([]byte(jsonAd), &adObj)
	if err != nil {
		return ad.Ad{}, err
	}

	return adObj, nil
}

func (s *Service) deleteFromRedis(ctx context.Context, ID int) {
	key := fmt.Sprintf("ad:%d", ID)
	s.rds.Rds.Del(ctx, key)
}
