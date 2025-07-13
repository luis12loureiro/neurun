package domain

import (
	"fmt"
)

type WorklowStatus string

const (
	WorkflowStatusIDLE      WorklowStatus = "IDLE"
	WorkflowStatusRunning   WorklowStatus = "RUNNING"
	WorkflowStatusCompleted WorklowStatus = "COMPLETED"
	WorkflowStatusFailed    WorklowStatus = "FAILED"
	// add more in the future...
)

type Worklow struct {
	ID          string
	Name        string
	Description string
	Status      WorklowStatus
	Tasks       []Task
}

type WorkflowRepository interface {
	Create(w Worklow) error
	Get(id string) (Worklow, error)
}

func (w *Worklow) String() string {
	return fmt.Sprintf("Id %s, Name %s", w.ID, w.Name)
}
