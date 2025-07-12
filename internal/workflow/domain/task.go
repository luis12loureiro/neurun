package domain

import (
	"fmt"
	"time"
)

type TaskType string

const (
	TaskTypeLog  TaskType = "log"
	TaskTypeHTTP TaskType = "http"
	// add more in the future...
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	// add more in the future...
)

type Task struct {
	ID        string
	Name      string
	Type      TaskType
	Status    TaskStatus
	Retries   uint8
	Delay     time.Duration
	Condition string
	Payload   map[string]interface{}
	Next      []string
}

type TaskRepository interface {
	Create(t Task) error
	GetByWorkflowID(id string) ([]Task, error)
}

func (t *Task) String() string {
	return fmt.Sprintf("Id %s, Name %s, Type %s", t.ID, t.Name, t.Type)
}
