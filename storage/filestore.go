package storage

import (
	"bulletin-board/domain"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type FileStore struct {
	filePath string
	mu       sync.Mutex
	ads      []domain.Ad
}

func (f *FileStore) NewBasePath(path string) {
	f.filePath = path
}

func (f *FileStore) List() ([]domain.Ad, error) {
	f.mu.Lock()
	defer f.mu.Lock()
	return f.readAll()
}

func (f *FileStore) Create(ad domain.Ad) (domain.Ad, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return domain.Ad{}, err
	}

	var maxId int64
	for _, item := range items {
		maxId = max(maxId, item.Id)
	}
	maxId++
	ad.Id = maxId

	items = append(items, ad)
	if err := f.writeAtomic(items); err != nil {
		return domain.Ad{}, err
	}
	return ad, nil
}

func (f *FileStore) Update(ad domain.Ad) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return err
	}

	updated := false
	for i := range items {
		if items[i].Id == ad.Id {
			updateItem(&items[i], &ad)
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("ad with id = %d not found", ad.Id)
	}

	return f.writeAtomic(items)
}

func (f *FileStore) Delete(id int64) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return err
	}

	newItems := make([]domain.Ad, len(items))

	deleted := false
	for _, item := range items {
		if item.Id == id {
			deleted = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !deleted {
		return fmt.Errorf("ad with id = %d not found", id)
	}

	return f.writeAtomic(newItems)
}

func (f *FileStore) writeAtomic(items []domain.Ad) error {
	dir := filepath.Dir(f.filePath)
	tmp, err := os.CreateTemp(dir, "tmp.json")
	if err != nil {
		return err
	}

	defer func() {
		_ = os.Remove(tmp.Name())
	}()

	enc := json.NewEncoder(tmp)
	enc.SetIndent("", "  ")
	if err := enc.Encode(items); err != nil {
		return err
	}

	if err := tmp.Sync(); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	return os.Rename(tmp.Name(), f.filePath)
}

func (f *FileStore) readAll() ([]domain.Ad, error) {
	file, err := os.Open(f.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return []domain.Ad{}, nil
	}
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var items []domain.Ad

	dec := json.NewDecoder(file)
	if err := dec.Decode(&items); err != nil {
		if errors.Is(err, io.EOF) {
			return []domain.Ad{}, nil
		}
		return nil, err
	}
	return items, nil
}

func updateItem(oldItem, newItem *domain.Ad) {
	oldItem.Title = newItem.Title
	oldItem.Description = newItem.Description
	oldItem.Price = newItem.Price
	oldItem.Contact = newItem.Contact
}
