package task

import (
	d "github.com/luis12loureiro/neurun/internal/task/domain"
)

type Repository interface {
	Create(t d.Task) error
	GetTask(id string) (d.Task, error)
}

type Service interface {
	Create(t d.Task) error
	GetTask(id string) (d.Task, error)
}

type service struct {
	r Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Create(t d.Task) error {
	return s.r.Create(t)
}

func (s *service) GetTask(id string) (d.Task, error) {
	return s.r.GetTask(id)
}
