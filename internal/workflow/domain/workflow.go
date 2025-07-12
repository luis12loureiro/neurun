package domain

import (
	"fmt"
	"time"
)

type Worklow struct {
	ID          string
	Name        string
	Description string
	Tasks       []Task
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type WorkflowRepository interface {
	Create(w Worklow) error
	Get(id string) (Worklow, error)
}

func (w *Worklow) String() string {
	return fmt.Sprintf("Id %s, Name %s", w.ID, w.Name)
}
