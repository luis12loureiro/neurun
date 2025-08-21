package domain

import (
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

type TaskType string

const (
	TaskTypeUnspecified TaskType = "UNSPECIFIED"
	TaskTypeLog         TaskType = "LOG"
	TaskTypeHTTP        TaskType = "HTTP"
	// add more in the future...
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "PENDING"
	TaskStatusRunning   TaskStatus = "RUNNING"
	TaskStatusCompleted TaskStatus = "COMPLETED"
	TaskStatusFailed    TaskStatus = "FAILED"
	// add more in the future...
)

type Task struct {
	ID         string
	Name       string
	Type       TaskType
	Status     TaskStatus
	Retries    uint8
	RetryDelay time.Duration
	Condition  string
	Payload    map[string]interface{}
	Next       []*Task
}

type TaskRepository interface {
	Create(t Task) error
	GetByWorkflowID(id string) ([]Task, error)
}

func NewTask(
	name string,
	taskType TaskType,
	retries uint32,
	retryDelay time.Duration,
	condition string,
	payload map[string]interface{},
	next []*Task,
) (*Task, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	// retries is uint8, so we need to check if it's larger than 255
	// the function receives uint32 on purpose
	if retries > math.MaxUint8 {
		return nil, fmt.Errorf("retries too large to fit in uint8")
	}
	if retryDelay < 0 {
		return nil, fmt.Errorf("retry delay cannot be negative")
	}
	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}
	return &Task{
		ID:         uuid.NewString(),
		Name:       name,
		Type:       taskType,
		Status:     TaskStatusPending,
		Retries:    uint8(retries),
		RetryDelay: retryDelay,
		Condition:  condition,
		Payload:    payload,
		Next:       next,
	}, nil
}

func (t *Task) String() string {
	return fmt.Sprintf("Id %s, Name %s, Type %s", t.ID, t.Name, t.Type)
}
