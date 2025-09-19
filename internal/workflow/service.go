package workflow

import (
	"context"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type Service interface {
	Create(w *domain.Workflow) (*domain.Workflow, error)
	Get(id string) (*domain.Workflow, error)
	Execute(ctx context.Context, id string) error
}

type service struct {
	r  domain.WorkflowRepository
	we WorkflowExecutor
}

func NewService(r domain.WorkflowRepository, we WorkflowExecutor) Service {
	return &service{
		r:  r,
		we: we,
	}
}

func (s *service) Create(w *domain.Workflow) (*domain.Workflow, error) {
	return w, s.r.Create(w)
}

func (s *service) Get(id string) (*domain.Workflow, error) {
	return s.r.Get(id)
}

func (s *service) Execute(ctx context.Context, id string) error {
	w, err := s.r.Get(id)
	if err != nil {
		return err
	}
	return s.we.Execute(ctx, w)
}
