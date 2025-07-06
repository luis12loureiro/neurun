package task

import (
	"fmt"
	"time"
)

type TaskType string

const (
	TaskTypeLog  TaskType = "log"
	TaskTypeHTTP TaskType = "http"
	// add more in the future...
)

type Task struct {
	Id        string
	Name      string
	Type      TaskType
	Retries   uint8
	Delay     time.Duration
	Condition string
	Payload   map[string]interface{}
	Next      []string
}

type ResultStatus string

const (
	Ok    ResultStatus = "ok"
	Error ResultStatus = "error"
)

type TaskResult struct {
	TaskId  string
	Status  ResultStatus
	Output  map[string][]interface{}
	Message string
}

func (t *Task) Execute() (*TaskResult, error) {
	if t == nil {
		return nil, fmt.Errorf("task does not exist")
	}

	if t.Delay > 0 {
		time.Sleep(t.Delay)
	}

	r := TaskResult{
		TaskId: t.Id,
		Status: Ok,
	}

	// TODO: define functions to execute each type
	switch t.Type {
	case TaskTypeLog:
		msg, ok := t.Payload["message"].(string)
		if !ok {
			r.Status = Error
			r.Message = "log task missing or invalid 'message' field"
		} else {
			fmt.Println("[LOG]", msg)
			r.Message = "log task successfully executed"
		}
	default:
		return nil, fmt.Errorf("unsupported task type: %s", t.Type)
	}

	return &r, nil
}

func (t *Task) String() string {
	return fmt.Sprintf("ID %s, Name %s, Type %s", t.Id, t.Name, t.Type)
}
