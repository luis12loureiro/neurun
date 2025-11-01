package workflow

import (
	"context"
	"sync"

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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	w.Status = domain.WorkflowStatusRunning
	// TODO: persist status change by having update method in repository
	for _, t := range w.Tasks {
		wg.Add(1)
		go func(task *domain.Task) {
			defer wg.Done()
			if err := we.executeTaskChain(ctx, task, resultCh); err != nil {
				select {
				case errCh <- err:
					cancel() // cancel all other tasks
				default: // error already sent, ignore
				}
			}
		}(t)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		w.Status = domain.WorkflowStatusFailed
		// TODO: persist status change by having update method in repository
		return err
	default:
		w.Status = domain.WorkflowStatusCompleted
		// TODO: persist status change by having update method in repository
		return nil
	}
}

func (we *workflowExecutor) executeTaskChain(ctx context.Context, task *domain.Task, resultCh chan<- map[string]interface{}) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	result, err := we.te.Execute(ctx, task)
	if err != nil {
		return err
	}

	resultCh <- map[string]interface{}{
		"task_id": task.ID,
		"status":  task.Status,
		"output":  result,
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	// execute next tasks if current task succeeded
	for _, nextTask := range task.Next {
		wg.Add(1)
		go func(nt *domain.Task) {
			defer wg.Done()
			if err := we.executeTaskChain(ctx, nt, resultCh); err != nil {
				select {
				case errCh <- err:
				default: // error already sent, ignore
				}
			}
		}(nextTask)
	}
	wg.Wait()

	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
