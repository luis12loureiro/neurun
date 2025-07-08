package domain

import (
	"fmt"
	"time"

	"github.com/luis12loureiro/neurun/internal/task/domain"
)

type Worklow struct {
	Id          string
	Name        string
	Description string
	Tasks       []domain.Task
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Repository interface {
	Create(w Worklow) error
	Get(id string) (Worklow, error)
}

func (w *Worklow) String() string {
	return fmt.Sprintf("Id %s, Name %s", w.Id, w.Name)
}
