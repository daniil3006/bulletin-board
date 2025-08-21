package dto

import "bulletin-board/internal/ad"

type RequestAd struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	UserID      int    `json:"user_id"`
}

type ResponseAd struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	UserID      int    `json:"user_id"`
}

func ToDto(ad ad.Ad) ResponseAd {
	return ResponseAd{
		ID:          ad.ID,
		Title:       ad.Title,
		Description: ad.Description,
		Price:       ad.Price,
		UserID:      ad.UserID,
	}
}

func ToAd(requestAd RequestAd) ad.Ad {
	return ad.Ad{
		Title:       requestAd.Title,
		Description: requestAd.Description,
		Price:       requestAd.Price,
		UserID:      requestAd.UserID,
	}
}
