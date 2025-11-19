package workflow

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
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

	// Build pending dependency counters for fan-in support
	pendingDeps := we.buildPendingDeps(w.Tasks)

	// Calculate total tasks in the workflow (count all tasks in pendingDeps map)
	totalTasks := len(pendingDeps)

	// track root tasks
	var wg sync.WaitGroup
	// buffered channel to capture first error without blocking
	errCh := make(chan error, 1)
	// track completed tasks to prevent cycles (thread-safe)
	completed := &sync.Map{}
	// track executed tasks count (thread-safe)
	executedCount := &atomic.Int32{}

	w.Status = domain.WorkflowStatusRunning
	for _, t := range w.Tasks {
		wg.Add(1) // increment wg counter
		go func(task *domain.Task) {
			defer wg.Done() // decrement wg counter
			if err := we.executeTaskChain(ctx, w, task, resultCh, pendingDeps, completed, executedCount, totalTasks); err != nil {
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
		resultCh <- map[string]interface{}{
			"taskId":         "",
			"status":         "",
			"output":         "",
			"workflowStatus": w.Status,
			"totalTasks":     totalTasks,
			"executedTasks":  int(executedCount.Load()),
		}
		return nil
	}
}

// buildPendingDeps traverses the task graph and builds atomic counters
// representing how many predecessors must complete before each task can run
func (we *workflowExecutor) buildPendingDeps(rootTasks []*domain.Task) map[string]*atomic.Int32 {
	pendingDeps := make(map[string]*atomic.Int32)
	visited := make(map[string]bool)

	var traverse func(task *domain.Task)
	traverse = func(task *domain.Task) {
		if visited[task.ID] {
			return
		}
		visited[task.ID] = true

		// initialize atomic counter to 0 if not exists
		if _, exists := pendingDeps[task.ID]; !exists {
			pendingDeps[task.ID] = &atomic.Int32{}
			pendingDeps[task.ID].Store(0)
		}

		// each next task has one more predecessor
		for _, nextTask := range task.Next {
			if _, exists := pendingDeps[nextTask.ID]; !exists {
				pendingDeps[nextTask.ID] = &atomic.Int32{}
				pendingDeps[nextTask.ID].Store(0)
			}
			pendingDeps[nextTask.ID].Add(1)
			traverse(nextTask)
		}
	}

	for _, root := range rootTasks {
		traverse(root)
	}

	return pendingDeps
}

func (we *workflowExecutor) executeTaskChain(
	ctx context.Context,
	w *domain.Workflow,
	task *domain.Task,
	resultCh chan<- map[string]interface{},
	pendingDeps map[string]*atomic.Int32,
	completed *sync.Map,
	executedCount *atomic.Int32,
	totalTasks int,
) error {
	// check for context cancellation
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// check for cycle: if task already completed, return error
	if _, alreadyCompleted := completed.Load(task.ID); alreadyCompleted {
		return fmt.Errorf("cycle detected: task %s (name: %s) already executed in this workflow execution", task.ID, task.Name)
	}

	// check if all dependencies are satisfied
	counter := pendingDeps[task.ID]
	if counter.Load() > 0 {
		// not ready yet, skip (another goroutine will execute when ready)
		return nil
	}

	// execute task
	result, err := we.te.Execute(ctx, task)
	if err != nil {
		return err
	}

	// mark task as completed
	completed.Store(task.ID, true)

	// increment executed count
	count := executedCount.Add(1)

	// stream result to channel
	resultCh <- map[string]interface{}{
		"taskId":         task.ID,
		"status":         task.Status,
		"output":         result,
		"workflowStatus": w.Status,
		"totalTasks":     totalTasks,
		"executedTasks":  int(count),
	}

	// track next tasks
	var wg sync.WaitGroup
	// buffered channel to capture first error without blocking
	errCh := make(chan error, 1)

	// execute next tasks if current task succeeded
	for _, nextTask := range task.Next {
		// atomically decrement the pending count for the next task
		newCount := pendingDeps[nextTask.ID].Add(-1)

		// if count reaches 0, all dependencies are satisfied
		if newCount == 0 {
			wg.Add(1) // increment wg counter
			go func(nt *domain.Task) {
				defer wg.Done() // decrement wg counter
				// recursively execute next tasks
				if err := we.executeTaskChain(ctx, w, nt, resultCh, pendingDeps, completed, executedCount, totalTasks); err != nil {
					select {
					case errCh <- err: // capture first error
					default: // error already sent, ignore
					}
				}
			}(nextTask)
		}
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
