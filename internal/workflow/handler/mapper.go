package handler

import (
	"github.com/luis12loureiro/neurun/internal/workflow/domain"

	pb "github.com/luis12loureiro/neurun/api/gen"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
)

func TaskFromProto(pbTask *pb.CreateTaskRequest) (*domain.Task, error) {
	next, err := convertNextFromProto(pbTask.GetNext())
	if err != nil {
		return nil, err
	}
	task, err := domain.NewTask(
		pbTask.GetName(),
		convertTaskTypeFromProto(pbTask.GetType()),
		pbTask.GetRetries(),
		pbTask.GetRetryDelay().AsDuration(),
		pbTask.GetCondition(),
		convertStructToMap(pbTask.GetPayload()),
		next,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func TaskToProto(t *domain.Task) *pb.Task {
	return &pb.Task{
		Id:         &t.ID,
		Name:       t.Name,
		Type:       convertTaskTypeToProto(t.Type),
		Status:     convertTaskStatusToProto(t.Status),
		Retries:    uint32(t.Retries),
		RetryDelay: durationpb.New(t.RetryDelay),
		Condition:  &t.Condition,
		Payload:    convertMapToStruct(t.Payload),
		Next:       convertNextToProto(t.Next),
	}
}

func convertTaskStatusToProto(s domain.TaskStatus) pb.TaskStatus {
	switch s {
	case domain.TaskStatusPending:
		return pb.TaskStatus_STATUS_PENDING
	case domain.TaskStatusRunning:
		return pb.TaskStatus_STATUS_RUNNNING
	case domain.TaskStatusCompleted:
		return pb.TaskStatus_STATUS_COMPLETED
	case domain.TaskStatusFailed:
		return pb.TaskStatus_STATUS_FAILED
	default:
		return pb.TaskStatus_STATUS_PENDING
	}
}

func convertStructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return nil
	}
	return s.AsMap()
}

func convertMapToStruct(m map[string]interface{}) *structpb.Struct {
	if m == nil {
		return nil
	}
	out, err := structpb.NewStruct(m)
	if err != nil {
		return nil
	}
	return out
}

func convertNextFromProto(pbNext []*pb.CreateTaskRequest) ([]*domain.Task, error) {
	if len(pbNext) == 0 {
		return []*domain.Task{}, nil
	}
	var out []*domain.Task
	for _, t := range pbNext {
		if t == nil || t.GetName() == "" {
			continue
		}
		from, err := TaskFromProto(t)
		if err != nil {
			return nil, err
		}
		out = append(out, from)
	}
	return out, nil
}

func convertNextToProto(t []*domain.Task) []*pb.Task {
	if len(t) == 0 {
		// protobuf reads nil as empty slice and retuns an emtpy array
		return nil
	}
	out := make([]*pb.Task, len(t))
	for i, t := range t {
		out[i] = TaskToProto(t)
	}
	return out
}

func convertTaskTypeFromProto(tt pb.TaskType) domain.TaskType {
	switch tt {
	case pb.TaskType_UNSPECIFIED:
		return domain.TaskTypeUnspecified
	case pb.TaskType_LOG:
		return domain.TaskTypeLog
	case pb.TaskType_HTTP:
		return domain.TaskTypeHTTP
	default:
		return domain.TaskTypeUnspecified
	}
}

func convertTaskTypeToProto(tt domain.TaskType) pb.TaskType {
	switch tt {
	case domain.TaskTypeUnspecified:
		return pb.TaskType_UNSPECIFIED
	case domain.TaskTypeLog:
		return pb.TaskType_LOG
	case domain.TaskTypeHTTP:
		return pb.TaskType_HTTP
	default:
		return pb.TaskType_UNSPECIFIED
	}
}
