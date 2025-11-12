package repository

import (
	"context"
	"sync"
)

type Repository interface {
	SaveProcessing(ctx context.Context, links []string) (int, error)
	SaveStorage(linksID int, res map[string]string)
	Get(ctx context.Context, linksID int) ([]string, map[string]string, error)
	GetIds(ctx context.Context) ([]int, error)
}

type repository struct {
	mu         sync.Mutex
	storage    map[int]map[string]string
	curId      int
	processing map[int][]string
	done       map[int][]string
}

func NewRepository() Repository {
	return &repository{
		storage:    map[int]map[string]string{},
		processing: make(map[int][]string),
		done:       make(map[int][]string),
		curId:      1,
	}
}

func (r *repository) SaveProcessing(ctx context.Context, links []string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id := r.curId
	r.curId++
	r.processing[id] = links
	return id, nil
}

func (r *repository) SaveStorage(linksID int, res map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.done[linksID] = r.processing[linksID]
	delete(r.processing, linksID)
	r.storage[linksID] = res
}

func (r *repository) Get(ctx context.Context, linksID int) ([]string, map[string]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if links, ok := r.done[linksID]; ok {
		return links, r.storage[linksID], nil
	}
	return nil, r.storage[linksID], nil
}

func (r *repository) GetIds(ctx context.Context) ([]int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var ids []int
	for id := range r.storage {
		ids = append(ids, id)
	}
	return ids, nil
}
