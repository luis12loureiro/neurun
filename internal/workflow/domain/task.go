package domain

import (
	"fmt"
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

const (
	TaskNameMaxLength        = 30
	TaskDescriptionMaxLength = 100
	TaskMaxRetries           = 5
	TaskMaxRetryDelay        = 5 * time.Second
	TaskMinRetryDelay        = 0 * time.Second
	TaskMaxNextLength        = 3
)

type Task struct {
	ID         string
	Name       string
	Type       TaskType
	Status     TaskStatus
	Retries    uint8
	RetryDelay time.Duration
	Condition  string
	Payload    Payload
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
	payload Payload,
	next []*Task,
) (*Task, error) {
	if name == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	if len([]rune(name)) > TaskNameMaxLength {
		return nil, fmt.Errorf("name cannot be longer than %d characters", TaskNameMaxLength)
	}
	if taskType == "" || taskType == TaskTypeUnspecified {
		return nil, fmt.Errorf("task type cannot be empty")
	}
	if retries > TaskMaxRetries {
		return nil, fmt.Errorf("retries cannot be more than %d", TaskMaxRetries)
	}
	if retryDelay < TaskMinRetryDelay || retryDelay > TaskMaxRetryDelay {
		return nil, fmt.Errorf("retry delay must be between %v and %v", TaskMinRetryDelay, TaskMaxRetryDelay)
	}
	if payload == nil {
		return nil, fmt.Errorf("payload cannot be nil")
	}
	if payload.Type() != taskType {
		return nil, fmt.Errorf("payload type does not match task type")
	}
	if len(next) > TaskMaxNextLength {
		return nil, fmt.Errorf("cannot have more than %d next tasks", TaskMaxNextLength)
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
