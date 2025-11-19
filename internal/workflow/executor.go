package workflow

import (
	"context"
	"fmt"
	"sync"
	"time"

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
	// create timeout context and defer cancel to cleanup
	ctx, timeout := context.WithTimeout(ctx, 5*time.Minute)
	defer timeout()
	// create cancelable context and defer cancel to cleanup
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// track root tasks
	var wg sync.WaitGroup
	// buffered channel to capture first error without blocking
	errCh := make(chan error, 1)
	// track visited tasks to prevent cycles (thread-safe)
	visited := &sync.Map{}

	w.Status = domain.WorkflowStatusRunning
	for _, t := range w.Tasks {
		wg.Add(1) // increment wg counter
		go func(task *domain.Task) {
			defer wg.Done() // decrement wg counter
			if err := we.executeTaskChain(ctx, task, resultCh, visited); err != nil {
				select {
				case errCh <- err:
					cancel() // cancel all other tasks
				default: // error already sent, ignore
				}
			}
		}(t)
	}
	// block until all tasks are done
	wg.Wait()

	select {
	case err := <-errCh:
		w.Status = domain.WorkflowStatusFailed
		return err
	default:
		w.Status = domain.WorkflowStatusCompleted
		return nil
	}
}

func (we *workflowExecutor) executeTaskChain(ctx context.Context, task *domain.Task, resultCh chan<- map[string]interface{}, visited *sync.Map) error {
	// check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// check for cycle: if task already visited, return error
	if _, alreadyVisited := visited.LoadOrStore(task.ID, true); alreadyVisited {
		return fmt.Errorf("cycle detected: task %s (name: %s) already executed in this workflow execution", task.ID, task.Name)
	}

	// execute task
	result, err := we.te.Execute(ctx, task)
	if err != nil {
		return err
	}

	// stream result to channel
	resultCh <- map[string]interface{}{
		"task_id": task.ID,
		"status":  task.Status,
		"output":  result,
	}

	// track next tasks
	var wg sync.WaitGroup
	// buffered channel to capture first error without blocking
	errCh := make(chan error, 1)

	// execute next tasks if current task succeeded
	for _, nextTask := range task.Next {
		wg.Add(1) // increment wg counter
		go func(nt *domain.Task) {
			defer wg.Done() // decrement wg counter
			// recursively execute next tasks
			if err := we.executeTaskChain(ctx, nt, resultCh, visited); err != nil {
				select {
				case errCh <- err: // capture first error
				default: // error already sent, ignore
				}
			}
		}(nextTask)
	}
	// block until all sibling tasks are done
	wg.Wait()

	// return first error if any
	select {
	case err := <-errCh:
		return err
	default:
		return nil
	}
}
