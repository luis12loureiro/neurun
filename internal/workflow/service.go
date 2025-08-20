package workflow

import (
	"github.com/google/uuid"
	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type Service interface {
	Create(w domain.Worklow) (domain.Worklow, error)
	Get(id string) (domain.Worklow, error)
}

type service struct {
	r domain.WorkflowRepository
}

func NewService(r domain.WorkflowRepository) Service {
	return &service{r}
}

func (s *service) Create(w domain.Worklow) (domain.Worklow, error) {
	w.ID = uuid.NewString()
	w.Status = domain.WorkflowStatusIDLE
	for i := range w.Tasks {
		w.Tasks[i].ID = uuid.NewString()
		w.Tasks[i].Status = domain.TaskStatusPending
	}
	return w, s.r.Create(w)
}

func (s *service) Get(id string) (domain.Worklow, error) {
	return s.r.Get(id)
}
