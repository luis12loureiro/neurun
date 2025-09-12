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

const (
	WorkflowNameMaxLength        = 30
	WorkflowDescriptionMaxLength = 100
	WorkflowMaxTasks             = 10
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
	if len([]rune(name)) > WorkflowNameMaxLength {
		return nil, fmt.Errorf("name cannot be longer than %d characters", WorkflowNameMaxLength)
	}
	if len([]rune(description)) > WorkflowDescriptionMaxLength {
		return nil, fmt.Errorf("description cannot be longer than %d characters", WorkflowDescriptionMaxLength)
	}
	totalTasks := countAllTasks(tasks)
	if totalTasks > WorkflowMaxTasks {
		return nil, fmt.Errorf("cannot have more than %d total tasks", WorkflowMaxTasks)
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

func countAllTasks(tasks []*Task) int {
	if len(tasks) == 0 {
		return 0
	}
	visited := make(map[string]bool)
	return countTasksRecursive(tasks, visited)
}

func countTasksRecursive(tasks []*Task, visited map[string]bool) int {
	count := 0
	for _, task := range tasks {
		if !visited[task.ID] {
			visited[task.ID] = true
			count++
			count += countTasksRecursive(task.Next, visited)
		}
	}
	return count
}
