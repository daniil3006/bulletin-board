package storage

import "bulletin-board/domain"

type Store interface {
	List() ([]domain.Ad, error)
	Create(ad domain.Ad) (domain.Ad, error)
	Update(ad domain.Ad) error
	Delete(id int) error
}
