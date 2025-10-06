package services

import (
	"sync"

	"github.com/afrizalsebastian/ai-cv-evaluator-with-go/models"
)

type IJobStore interface {
	Set(id string, job *models.JobItem)
	Get(id string) (*models.JobItem, bool)
	Delete(id string)
}

type jobStore struct {
	jobs map[string]*models.JobItem
	mu   sync.RWMutex
}

func NewJobStore() IJobStore {
	return &jobStore{
		jobs: make(map[string]*models.JobItem),
	}
}

func (js *jobStore) Set(id string, job *models.JobItem) {
	js.mu.Lock()
	defer js.mu.Unlock()
	js.jobs[id] = job
}

func (js *jobStore) Get(id string) (*models.JobItem, bool) {
	js.mu.RLock()
	defer js.mu.RUnlock()
	job, ok := js.jobs[id]
	return job, ok
}

func (js *jobStore) Delete(id string) {
	js.mu.Lock()
	defer js.mu.Unlock()
	delete(js.jobs, id)
}
