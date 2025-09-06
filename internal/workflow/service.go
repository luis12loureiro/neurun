package workflow

import (
	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type Service interface {
	Create(w *domain.Workflow) (*domain.Workflow, error)
	Get(id string) (*domain.Workflow, error)
}

type service struct {
	r domain.WorkflowRepository
}

func NewService(r domain.WorkflowRepository) Service {
	return &service{r}
}

func (s *service) Create(w *domain.Workflow) (*domain.Workflow, error) {
	return w, s.r.Create(w)
}

func (s *service) Get(id string) (*domain.Workflow, error) {
	return s.r.Get(id)
}
