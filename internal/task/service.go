package task

import "github.com/luis12loureiro/neurun/internal/task/domain"

type Service interface {
	Create(t domain.Task) error
	Get(id string) (domain.Task, error)
}

type service struct {
	r domain.Repository
}

func NewService(r domain.Repository) Service {
	return &service{r}
}

func (s *service) Create(t domain.Task) error {
	return s.r.Create(t)
}

func (s *service) Get(id string) (domain.Task, error) {
	return s.r.Get(id)
}
