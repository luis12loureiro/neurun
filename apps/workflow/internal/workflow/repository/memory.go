package repository

import (
	"fmt"

	"github.com/luis12loureiro/neurun/apps/workflow/internal/workflow/domain"
	"github.com/luis12loureiro/neurun/apps/workflow/internal/workflow/repository/storage"
)

type MemoryRepo struct {
	workflows map[string]domain.Workflow
}

func NewMemoryRepository() domain.WorkflowRepository {
	return &MemoryRepo{
		workflows: storage.DefaultWorkflows,
	}
}

func (r *MemoryRepo) Create(t *domain.Workflow) error {
	r.workflows[t.ID] = *t
	return nil
}

func (r *MemoryRepo) Get(id string) (*domain.Workflow, error) {
	w, exists := r.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow with id %s not found", id)
	}
	return &w, nil
}
