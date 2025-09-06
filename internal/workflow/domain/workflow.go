package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type WorklowStatus string

const (
	WorkflowStatusIDLE      WorklowStatus = "IDLE"
	WorkflowStatusRunning   WorklowStatus = "RUNNING"
	WorkflowStatusCompleted WorklowStatus = "COMPLETED"
	WorkflowStatusFailed    WorklowStatus = "FAILED"
	// add more in the future...
)

type Workflow struct {
	ID          string
	Name        string
	Description string
	Status      WorklowStatus
	Tasks       []*Task
}

type WorkflowRepository interface {
	Create(w *Workflow) error
	Get(id string) (*Workflow, error)
}

func NewWorkflow(name string, description string, tasks []*Task) (*Workflow, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if len(tasks) > 10 {
		return nil, fmt.Errorf("cannot have more than 10 tasks")
	}
	return &Workflow{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Status:      WorkflowStatusIDLE,
		Tasks:       tasks,
	}, nil
}

func (w *Workflow) String() string {
	return fmt.Sprintf("Id %s, Name %s", w.ID, w.Name)
}
