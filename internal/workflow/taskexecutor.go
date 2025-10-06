package workflow

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/luis12loureiro/neurun/internal/workflow/domain"
)

type TaskExecutor interface {
	Execute(ctx context.Context, t *domain.Task) (interface{}, error)
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

func (te *taskExecutor) Execute(ctx context.Context, t *domain.Task) (interface{}, error) {
	t.Status = domain.TaskStatusRunning

	var err error
	var output interface{}
	switch t.Type {
	case domain.TaskTypeHTTP:
		return nil, fmt.Errorf("HTTP task executor not implemented yet")
	case domain.TaskTypeLog:
		logPayload, ok := t.Payload.(*domain.LogPayload)
		if !ok {
			err = fmt.Errorf("invalid payload type for LOG task")
		} else {
			randomSeconds := rand.Intn(3)
			time.Sleep(time.Duration(randomSeconds) * time.Second)
			output = logPayload.Message
		}
	default:
		err = fmt.Errorf("unknown task type: %v", t.Type)
	}

	if err != nil {
		t.Status = domain.TaskStatusFailed
		return nil, err
	} else {
		t.Status = domain.TaskStatusCompleted
	}
	return output, nil
}
