package task

import (
	"github.com/luis12loureiro/neurun/internal/task/domain"
)

type Service interface {
	Create(t domain.Task) error
	GetTask(id string) (domain.Task, error)
}

type service struct {
	r domain.Repository
}

func NewService(r domain.Repository) *service {
	return &service{r}
}

func (s *service) Create(t domain.Task) error {
	return s.r.Create(t)
}

func (s *service) GetTask(id string) (domain.Task, error) {
	return s.r.GetTask(id)
}
