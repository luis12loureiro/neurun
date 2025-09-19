package workflow

import (
	"context"
	"fmt"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type TaskExecutor interface {
	Execute(ctx context.Context, t *domain.Task) error
}

type taskExecutor struct {
	// httpExecutor HTTPTaskExecutor TODO: Add HTTPTaskExecutor
	// logExecutor  LogTaskExecutor TODO: Add LogTaskExecutor
}

func NewTaskExecutor() TaskExecutor {
	return &taskExecutor{
		// httpExecutor: NewHTTPTaskExecutor(), TODO: Add HTTPTaskExecutor
		// logExecutor:  NewLogTaskExecutor(),  TODO: Add LogTaskExecutor
	}
}

func (te *taskExecutor) Execute(ctx context.Context, t *domain.Task) error {
	t.Status = domain.TaskStatusRunning

	var err error
	switch t.Type {
	case domain.TaskTypeHTTP:
		return fmt.Errorf("HTTP task executor not implemented yet")
	case domain.TaskTypeLog:
		logPayload, ok := t.Payload.(*domain.LogPayload)
		if !ok {
			err = fmt.Errorf("invalid payload type for LOG task")
		} else {
			_, err = fmt.Println(logPayload.Message)
		}
	default:
		err = fmt.Errorf("unknown task type: %v", t.Type)
	}

	if err != nil {
		t.Status = domain.TaskStatusFailed
		return err
	} else {
		t.Status = domain.TaskStatusCompleted
	}
	return nil
}
