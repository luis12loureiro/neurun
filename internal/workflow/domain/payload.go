package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

type Payload interface {
	Type() TaskType
}

type LogPayload struct {
	Message string
}

func (l *LogPayload) Type() TaskType {
	return TaskTypeLog
}

type HTTPApiKeyLocation string

const (
	HTTPMinTimeout                              = 0 * time.Second
	HTTPMaxTimeout                              = 30 * time.Second
	HTTPMinStatusCode                           = 100
	HTTPMaxStatusCode                           = 599
	HTTPBasicAuthType                           = "basic"
	HTTPBearerAuthType                          = "bearer"
	HTTPApiKeyAuthType                          = "apikey"
	HTTPApiKeyLocationHeader HTTPApiKeyLocation = "HEADER"
	HTTPApiKeyLocationQuery  HTTPApiKeyLocation = "QUERY"
)

type HTTPPayload struct {
	URL                string
	Method             string
	Body               []byte
	Headers            map[string]string
	QueryParams        map[string]string
	Timeout            time.Duration // max time to wait for a response
	Auth               HTTPAuthType
	FollowRedirects    bool // whether to follow HTTP redirects
	VerifySSL          bool // whether to verify SSL certificates
	ExpectedStatusCode int32
}

func (h *HTTPPayload) Type() TaskType {
	return TaskTypeHTTP
}

func NewHTTPPayload(urlStr, method string, body []byte, headers, queryParams map[string]string,
	timeout time.Duration, auth HTTPAuthType, followRedirects, verifySSL bool,
	expectedStatusCode int32) (*HTTPPayload, error) {

	if urlStr == "" {
		return nil, fmt.Errorf("URL cannot be empty")
	}
	if _, err := url.Parse(urlStr); err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}
	if method == "" {
		return nil, fmt.Errorf("HTTP method cannot be empty")
	}
	validMethods := []string{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"}
	methodUpper := strings.ToUpper(method)
	valid := false
	for _, vm := range validMethods {
		if vm == methodUpper {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("invalid HTTP method: %s", method)
	}
	if timeout < HTTPMinTimeout || timeout > HTTPMaxTimeout {
		return nil, fmt.Errorf("timeout must be between %d and %d", HTTPMinTimeout, HTTPMaxTimeout)
	}
	if expectedStatusCode < HTTPMinStatusCode || expectedStatusCode > HTTPMaxStatusCode {
		return nil, fmt.Errorf("invalid HTTP status code: %d", expectedStatusCode)
	}
	if auth != nil {
		if err := validateHTTPAuth(auth); err != nil {
			return nil, fmt.Errorf("invalid auth: %w", err)
		}
	}
	if headers == nil {
		headers = make(map[string]string)
	}
	if queryParams == nil {
		queryParams = make(map[string]string)
	}

	return &HTTPPayload{
		URL:                urlStr,
		Method:             methodUpper,
		Body:               body,
		Headers:            headers,
		QueryParams:        queryParams,
		Timeout:            timeout,
		Auth:               auth,
		FollowRedirects:    followRedirects,
		VerifySSL:          verifySSL,
		ExpectedStatusCode: expectedStatusCode,
	}, nil
}

func validateHTTPAuth(auth HTTPAuthType) error {
	switch a := auth.(type) {
	case *HTTPBasicAuth:
		if a.Username == "" {
			return fmt.Errorf("basic auth username cannot be empty")
		}
		if a.Password == "" {
			return fmt.Errorf("basic auth password cannot be empty")
		}
	case *HTTPBearerAuth:
		if a.Token == "" {
			return fmt.Errorf("bearer token cannot be empty")
		}
	case *HTTPApiKeyAuth:
		if a.Key == "" {
			return fmt.Errorf("API key name cannot be empty")
		}
		if a.Value == "" {
			return fmt.Errorf("API key value cannot be empty")
		}
		if a.Location != HTTPApiKeyLocationHeader && a.Location != HTTPApiKeyLocationQuery {
			return fmt.Errorf("invalid API key location")
		}
	default:
		return fmt.Errorf("unknown auth type")
	}
	return nil
}

type HTTPAuthType interface {
	Type() string
}

type HTTPBasicAuth struct {
	Username string
	Password string
}

func (h *HTTPBasicAuth) Type() string {
	return HTTPBasicAuthType
}

type HTTPBearerAuth struct {
	Token string
}

func (h *HTTPBearerAuth) Type() string {
	return HTTPBearerAuthType
}

type HTTPApiKeyAuth struct {
	Key      string
	Value    string
	Location HTTPApiKeyLocation
}

func (h *HTTPApiKeyAuth) Type() string {
	return HTTPApiKeyAuthType
}
