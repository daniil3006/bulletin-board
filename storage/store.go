package storage

import "bulletin-board/domain"

type Store interface {
	GetAll() ([]domain.Ad, error)
	GetById(ID int64) (domain.Ad, error)
	Create(ad domain.Ad) (domain.Ad, error)
	Update(ad domain.Ad) (domain.Ad, error)
	Delete(id int64) error
}
