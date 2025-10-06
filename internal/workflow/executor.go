package workflow

import (
	"context"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type WorkflowExecutor interface {
	Execute(ctx context.Context, w *domain.Workflow, resultCh chan<- map[string]interface{}) error
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

func (we *workflowExecutor) Execute(ctx context.Context, w *domain.Workflow, resultCh chan<- map[string]interface{}) error {
	w.Status = domain.WorkflowStatusRunning
	// TODO: persist status change by having update method in repository
	for _, t := range w.Tasks {
		if err := we.executeTaskChain(ctx, t, resultCh); err != nil {
			w.Status = domain.WorkflowStatusFailed
			// TODO: persist status change by having update method in repository
			return err
		}
	}
	w.Status = domain.WorkflowStatusCompleted
	// TODO: persist status change by having update method in repository
	return nil
}

func (we *workflowExecutor) executeTaskChain(ctx context.Context, task *domain.Task, resultCh chan<- map[string]interface{}) error {
	// execute current task
	result, err := we.te.Execute(ctx, task)
	if err != nil {
		return err
	}

	resultCh <- map[string]interface{}{
		"task_id": task.ID,
		"status":  task.Status,
		"output":  result,
	}

	// execute next tasks if current task succeeded
	for _, nextTask := range task.Next {
		if err := we.executeTaskChain(ctx, nextTask, resultCh); err != nil {
			return err
		}
	}
	return nil
}
