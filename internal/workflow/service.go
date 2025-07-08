package workflow

import "github.com/luis12loureiro/neurun/internal/workflow/domain"

type Service interface {
	Create(w domain.Worklow) error
	Get(id string) (domain.Worklow, error)
}

type service struct {
	r domain.Repository
}

func NewService(r domain.Repository) Service {
	return &service{r}
}

func (s *service) Create(w domain.Worklow) error {
	return s.r.Create(w)
}

func (s *service) Get(id string) (domain.Worklow, error) {
	return s.r.Get(id)
}
