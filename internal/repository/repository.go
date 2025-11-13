package repository

import (
	"encoding/json"
	"os"
	"sync"
	"url-checker/internal/models"
)

type Repository interface {
	SaveTask(task *models.Task) error
	UpdateTask(t *models.Task) error
	GetTask(id int) (*models.Task, bool)
	GetTasks(ids []int) []*models.Task
	GetAllProcessingTasks() []*models.Task
}

type repository struct {
	mu       sync.RWMutex
	storage  map[int]*models.Task
	curId    int
	filePath string
}

func NewRepository(filePath string) (Repository, error) {
	r := &repository{
		storage:  make(map[int]*models.Task),
		curId:    0,
		filePath: filePath,
	}

	if err := r.loadFromFile(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *repository) SaveTask(t *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.curId++
	t.ID = r.curId
	r.storage[t.ID] = t

	return r.saveToFile()
}

func (r *repository) UpdateTask(t *models.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.storage[t.ID] = t
	return r.saveToFile()
}

func (r *repository) GetTask(id int) (*models.Task, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, exists := r.storage[id]
	return task, exists
}

func (r *repository) GetTasks(ids []int) []*models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*models.Task
	for _, id := range ids {
		if t, exists := r.storage[id]; exists {
			tasks = append(tasks, t)
		}
	}

	return tasks
}

func (r *repository) GetAllProcessingTasks() []*models.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tasks []*models.Task
	for _, t := range r.storage {
		if t.Status == "processing" {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (r *repository) saveToFile() error {
	data := struct {
		CurId   int                  `json:"cur_id"`
		Storage map[int]*models.Task `json:"storage"`
	}{
		CurId:   r.curId,
		Storage: r.storage,
	}

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(data)
}

func (r *repository) loadFromFile() error {
	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var data struct {
		CurId   int                  `json:"cur_id"`
		Storage map[int]*models.Task `json:"storage"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return err
	}

	r.curId = data.CurId
	r.storage = data.Storage
	if r.storage == nil {
		r.storage = make(map[int]*models.Task)
	}

	return nil
}
