package storage

import (
	"github.com/luis12loureiro/neurun/apps/workflow/internal/workflow/domain"
)

// Fan-In Test Workflow: 3 root tasks (A, B, C) all point to 1 shared task (D)
var (
	sharedTaskD = &domain.Task{
		ID:         "task-d",
		Name:       "Shared Task D",
		Type:       domain.TaskTypeLog,
		Status:     domain.TaskStatusPending,
		Retries:    0,
		RetryDelay: 0,
		Payload: &domain.LogPayload{
			Message: "✅ All 3 tasks completed! (A, B, C)",
		},
		Next: []*domain.Task{
			{
				ID:         "task-e",
				Name:       "Task E",
				Type:       domain.TaskTypeLog,
				Status:     domain.TaskStatusPending,
				Retries:    0,
				RetryDelay: 0,
				Payload: &domain.LogPayload{
					Message: "Task E executing",
				},
			},
		},
	}

	DefaultWorkflows = map[string]domain.Workflow{
		"fan-in-test": {
			ID:          "fan-in-test",
			Name:        "Fan-In Test",
			Description: "3 root tasks → 1 shared task",
			Status:      domain.WorkflowStatusIDLE,
			Tasks: []*domain.Task{
				{
					ID:         "task-a",
					Name:       "Root Task A",
					Type:       domain.TaskTypeLog,
					Status:     domain.TaskStatusPending,
					Retries:    0,
					RetryDelay: 0,
					Payload: &domain.LogPayload{
						Message: "🅰️  Task A executing",
					},
					Next: []*domain.Task{sharedTaskD},
				},
				{
					ID:         "task-b",
					Name:       "Root Task B",
					Type:       domain.TaskTypeLog,
					Status:     domain.TaskStatusPending,
					Retries:    0,
					RetryDelay: 0,
					Payload: &domain.LogPayload{
						Message: "🅱️  Task B executing",
					},
					Next: []*domain.Task{sharedTaskD},
				},
				{
					ID:         "task-c",
					Name:       "Root Task C",
					Type:       domain.TaskTypeLog,
					Status:     domain.TaskStatusPending,
					Retries:    0,
					RetryDelay: 0,
					Payload: &domain.LogPayload{
						Message: "🅲  Task C executing",
					},
					Next: []*domain.Task{sharedTaskD},
				},
			},
		},
	}
)
