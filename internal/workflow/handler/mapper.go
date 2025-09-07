package handler

import (
	"github.com/luis12loureiro/neurun/internal/workflow/domain"

	pb "github.com/luis12loureiro/neurun/api/gen"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TaskFromProto(pbTask *pb.CreateTaskRequest) (*domain.Task, error) {
	next, err := convertNextFromProto(pbTask.GetNext())
	if err != nil {
		return nil, err
	}
	var payload domain.Payload
	switch pbTask.GetPayload().(type) {
	case *pb.CreateTaskRequest_LogPayload:
		payload = &domain.LogPayload{
			Message: pbTask.GetLogPayload().GetMessage(),
		}
	case *pb.CreateTaskRequest_HttpPayload:
		httpPayload := pbTask.GetHttpPayload()
		httpPayloadDomain, err := domain.NewHTTPPayload(
			httpPayload.GetUrl(),
			httpPayload.GetMethod(),
			httpPayload.GetBody(),
			httpPayload.GetHeaders(),
			httpPayload.GetQueryParams(),
			httpPayload.GetTimeout().AsDuration(),
			convertHTTPAuthFromProto(httpPayload.GetAuth()),
			httpPayload.GetFollowRedirects(),
			httpPayload.GetVerifySSL(),
			httpPayload.GetExpectedStatusCode(),
		)
		if err != nil {
			return nil, err
		}
		payload = httpPayloadDomain
	}
	task, err := domain.NewTask(
		pbTask.GetName(),
		convertTaskTypeFromProto(pbTask.GetType()),
		pbTask.GetRetries(),
		pbTask.GetRetryDelay().AsDuration(),
		pbTask.GetCondition(),
		payload,
		next,
	)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func TaskToProto(t *domain.Task) *pb.Task {
	switch p := t.Payload.(type) {
	case *domain.LogPayload:
		return &pb.Task{
			Id:         &t.ID,
			Name:       t.Name,
			Type:       convertTaskTypeToProto(t.Type),
			Status:     convertTaskStatusToProto(t.Status),
			Retries:    uint32(t.Retries),
			RetryDelay: durationpb.New(t.RetryDelay),
			Condition:  &t.Condition,
			Payload: &pb.Task_LogPayload{
				LogPayload: &pb.LogPayload{
					Message: p.Message,
				},
			},
			Next: convertNextToProto(t.Next),
		}
	case *domain.HTTPPayload:
		return &pb.Task{
			Id:         &t.ID,
			Name:       t.Name,
			Type:       convertTaskTypeToProto(t.Type),
			Status:     convertTaskStatusToProto(t.Status),
			Retries:    uint32(t.Retries),
			RetryDelay: durationpb.New(t.RetryDelay),
			Condition:  &t.Condition,
			Payload: &pb.Task_HttpPayload{
				HttpPayload: &pb.HTTPPayload{
					Url:                p.URL,
					Method:             p.Method,
					Body:               p.Body,
					Headers:            p.Headers,
					QueryParams:        p.QueryParams,
					Timeout:            durationpb.New(p.Timeout),
					Auth:               convertHTTPAuthToProto(p.Auth),
					FollowRedirects:    p.FollowRedirects,
					VerifySSL:          p.VerifySSL,
					ExpectedStatusCode: p.ExpectedStatusCode,
				},
			},
			Next: convertNextToProto(t.Next),
		}
	default:
		return &pb.Task{
			Id:         &t.ID,
			Name:       t.Name,
			Type:       convertTaskTypeToProto(t.Type),
			Status:     convertTaskStatusToProto(t.Status),
			Retries:    uint32(t.Retries),
			RetryDelay: durationpb.New(t.RetryDelay),
			Condition:  &t.Condition,
			Next:       convertNextToProto(t.Next),
		}
	}
}

func convertTaskStatusToProto(s domain.TaskStatus) pb.TaskStatus {
	switch s {
	case domain.TaskStatusPending:
		return pb.TaskStatus_STATUS_PENDING
	case domain.TaskStatusRunning:
		return pb.TaskStatus_STATUS_RUNNING
	case domain.TaskStatusCompleted:
		return pb.TaskStatus_STATUS_COMPLETED
	case domain.TaskStatusFailed:
		return pb.TaskStatus_STATUS_FAILED
	default:
		return pb.TaskStatus_STATUS_PENDING
	}
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

func convertHTTPAuthFromProto(pbAuth *pb.HTTPAuth) domain.HTTPAuthType {
	if pbAuth == nil {
		return nil
	}

	switch authType := pbAuth.GetAuthType().(type) {
	case *pb.HTTPAuth_Basic:
		return &domain.HTTPBasicAuth{
			Username: authType.Basic.GetUsername(),
			Password: authType.Basic.GetPassword(),
		}
	case *pb.HTTPAuth_Bearer:
		return &domain.HTTPBearerAuth{
			Token: authType.Bearer.GetToken(),
		}
	case *pb.HTTPAuth_ApiKey:
		return &domain.HTTPApiKeyAuth{
			Key:      authType.ApiKey.GetKey(),
			Value:    authType.ApiKey.GetValue(),
			Location: convertHTTPApiKeyLocationFromProto(authType.ApiKey.GetLocation()),
		}
	default:
		return nil
	}
}

func convertHTTPAuthToProto(auth domain.HTTPAuthType) *pb.HTTPAuth {
	if auth == nil {
		return nil
	}

	pbAuth := &pb.HTTPAuth{}
	switch authType := auth.(type) {
	case *domain.HTTPBasicAuth:
		pbAuth.AuthType = &pb.HTTPAuth_Basic{
			Basic: &pb.HTTPBasicAuth{
				Username: authType.Username,
				Password: authType.Password,
			},
		}
	case *domain.HTTPBearerAuth:
		pbAuth.AuthType = &pb.HTTPAuth_Bearer{
			Bearer: &pb.HTTPBearerAuth{
				Token: authType.Token,
			},
		}
	case *domain.HTTPApiKeyAuth:
		pbAuth.AuthType = &pb.HTTPAuth_ApiKey{
			ApiKey: &pb.HTTPApiKeyAuth{
				Key:      authType.Key,
				Value:    authType.Value,
				Location: convertHTTPApiKeyLocationToProto(authType.Location),
			},
		}
	}
	return pbAuth
}

func convertHTTPApiKeyLocationFromProto(location pb.HTTPApiKeyLocation) domain.HTTPApiKeyLocation {
	switch location {
	case pb.HTTPApiKeyLocation_HEADER:
		return domain.HTTPApiKeyLocationHeader
	case pb.HTTPApiKeyLocation_QUERY:
		return domain.HTTPApiKeyLocationQuery
	default:
		return domain.HTTPApiKeyLocationHeader
	}
}

func convertHTTPApiKeyLocationToProto(location domain.HTTPApiKeyLocation) pb.HTTPApiKeyLocation {
	switch location {
	case domain.HTTPApiKeyLocationHeader:
		return pb.HTTPApiKeyLocation_HEADER
	case domain.HTTPApiKeyLocationQuery:
		return pb.HTTPApiKeyLocation_QUERY
	default:
		return pb.HTTPApiKeyLocation_HEADER
	}
}
