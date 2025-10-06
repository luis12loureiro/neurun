package storage

import (
	"time"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

var DefaultWorkflows = map[string]domain.Workflow{
	"wf1": {
		ID:          "wf1",
		Name:        "Default Workflow 1",
		Description: "This is the first default workflow",
		Status:      domain.WorkflowStatusIDLE,
		Tasks: []*domain.Task{
			{
				ID:         "task1",
				Name:       "Task 1",
				Type:       domain.TaskTypeLog,
				Status:     domain.TaskStatusPending,
				Retries:    1,
				RetryDelay: 2 * time.Second,
				Payload: &domain.LogPayload{
					Message: "Hello, World!",
				},
				Next: []*domain.Task{
					{
						ID:         "task2",
						Name:       "Task 2",
						Type:       domain.TaskTypeLog,
						Status:     domain.TaskStatusPending,
						Retries:    2,
						RetryDelay: 3 * time.Second,
						Payload: &domain.LogPayload{
							Message: "My name is Luís!",
						},
					},
				},
			},
			{
				ID:         "task3",
				Name:       "Task 3",
				Type:       domain.TaskTypeLog,
				Status:     domain.TaskStatusPending,
				Retries:    1,
				RetryDelay: 2 * time.Second,
				Payload: &domain.LogPayload{
					Message: "Hallo, Welt!",
				},
				Next: []*domain.Task{
					{
						ID:         "task4",
						Name:       "Task 4",
						Type:       domain.TaskTypeLog,
						Status:     domain.TaskStatusPending,
						Retries:    2,
						RetryDelay: 3 * time.Second,
						Payload: &domain.LogPayload{
							Message: "Ich bin Luís!",
						},
					},
				},
			},
		},
	},
	"wf2": {
		ID:          "wf2",
		Name:        "Default Workflow 2",
		Description: "This is the second default workflow",
		Status:      domain.WorkflowStatusIDLE,
		Tasks:       []*domain.Task{},
	},
}
