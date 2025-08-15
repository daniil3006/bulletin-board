package fileStore

import (
	"bulletin-board/internal/ad"
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type fileStore struct {
	filePath string
	mu       sync.Mutex
}

func (f fileStore) GetAll(ctx context.Context) ([]ad.Ad, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.readAll()
}

func (f fileStore) GetByID(ctx context.Context, ID int) (ad.Ad, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	items, err := f.readAll()
	if err != nil {
		return ad.Ad{}, err
	}

	for _, item := range items {
		if item.ID == ID {
			return item, nil
		}
	}
	return ad.Ad{}, ad.ErrNotFound
}

func (f fileStore) Create(ctx context.Context, newAd ad.Ad) (ad.Ad, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return ad.Ad{}, err
	}

	var maxId int
	for _, item := range items {
		maxId = max(maxId, item.ID)
	}
	maxId++
	newAd.ID = maxId

	items = append(items, newAd)
	if err := f.writeAtomic(items); err != nil {
		return ad.Ad{}, err
	}
	return newAd, nil
}

func (f fileStore) Update(ctx context.Context, newAd ad.Ad, id int) (ad.Ad, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return ad.Ad{}, err
	}

	updated := false
	for i := range items {
		if items[i].ID == newAd.ID {
			updateItem(&items[i], &newAd)
			updated = true
			break
		}
	}

	if !updated {
		return ad.Ad{}, ad.ErrNotFound
	}

	return newAd, f.writeAtomic(items)
}

func (f fileStore) Delete(ctx context.Context, ID int) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	items, err := f.readAll()
	if err != nil {
		return err
	}

	newItems := make([]ad.Ad, 0, len(items))

	deleted := false
	for _, item := range items {
		if item.ID == ID {
			deleted = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !deleted {
		return ad.ErrNotFound
	}

	return f.writeAtomic(newItems)
}

func NewRepository(path string) ad.Repository {
	return &fileStore{filePath: path}
}

func (f *fileStore) writeAtomic(items []ad.Ad) error {
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

func (f *fileStore) readAll() ([]ad.Ad, error) {
	file, err := os.Open(f.filePath)
	if errors.Is(err, os.ErrNotExist) {
		return []ad.Ad{}, nil
	}
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var items []ad.Ad

	dec := json.NewDecoder(file)
	if err := dec.Decode(&items); err != nil {
		if errors.Is(err, io.EOF) {
			return []ad.Ad{}, nil
		}
		return nil, err
	}
	return items, nil
}

func updateItem(oldItem, newItem *ad.Ad) {
	oldItem.Title = newItem.Title
	oldItem.Description = newItem.Description
	oldItem.Price = newItem.Price
	oldItem.Contact = newItem.Contact
}
