package storage

import (
	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

var DefaultWorkflows = map[string]domain.Workflow{
	"wf1": {
		ID:          "wf1",
		Name:        "Default Workflow 1",
		Description: "This is the first default workflow",
		Status:      domain.WorkflowStatusIDLE,
		Tasks:       []*domain.Task{},
	},
	"wf2": {
		ID:          "wf2",
		Name:        "Default Workflow 2",
		Description: "This is the second default workflow",
		Status:      domain.WorkflowStatusIDLE,
		Tasks:       []*domain.Task{},
	},
}
