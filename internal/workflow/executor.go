package workflow

import (
	"context"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type WorkflowExecutor interface {
	Execute(ctx context.Context, w *domain.Workflow) error
}

type workflowExecutor struct {
	te TaskExecutor
	r  domain.WorkflowRepository
}

func NewWorkflowExecutor(r domain.WorkflowRepository, te TaskExecutor) WorkflowExecutor {
	return &workflowExecutor{
		r:  r,
		te: te,
	}
}

func (we *workflowExecutor) Execute(ctx context.Context, w *domain.Workflow) error {
	w.Status = domain.WorkflowStatusRunning
	// TODO: persist status change by having update method in repository
	for _, t := range w.Tasks {
		if err := we.executeTaskChain(ctx, t); err != nil {
			w.Status = domain.WorkflowStatusFailed
			// TODO: persist status change by having update method in repository
			return err
		}
	}
	w.Status = domain.WorkflowStatusCompleted
	// TODO: persist status change by having update method in repository
	return nil
}

func (we *workflowExecutor) executeTaskChain(ctx context.Context, task *domain.Task) error {
	// execute current task
	if err := we.te.Execute(ctx, task); err != nil {
		return err
	}
	// execute next tasks if current task succeeded
	for _, nextTask := range task.Next {
		if err := we.executeTaskChain(ctx, nextTask); err != nil {
			return err
		}
	}
	return nil
}
