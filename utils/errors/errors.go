package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	//Persistence errors
	QueryError    = "db_query_error"
	NotFoundError = "db_not_found_error"
	CreateError   = "db_create_error"
	SaveError     = "db_save_error"
	UpdateError   = "db_update_error"
	DeleteError   = "db_delete_error"

	//Proxy errors
	ProxyRequestError         = "proxy_request_error"
	ProxyInvalidResponseError = "proxy_invalid_response_error"

	//Domain errors
	NotValid = "entity_not_valid"

	//Marshaller errors
	UnmarshalError = "unmarshal_error"

	// Pipeline errors
	PipelineProcessingError = "pipeline_processing_error"
	CleanupError = "pipeline_cleanup_error"

	// Handler errors
	ServerInvalidParamsError = "server_handler_invalid_params_error"
	CliInvalidParamsError    = "cli_handler_invalid_params_error"
	JobInvalidParamsError    = "job_handler_invalid_params_error"
)

type ApplicationError interface {
	Message() string
	Status() string
	Error() string
}

type applicationError struct {
	ErrMessage string `json:"message"`
	ErrStatus  string `json:"status"`
	ErrError   error  `json:"error"`
}

func (e applicationError) Error() string {
	return fmt.Sprintf("message: %s - status: %s - error: %s",
		e.ErrMessage, e.ErrStatus, e.ErrError)
}

func (e applicationError) Message() string {
	return e.ErrMessage
}

func (e applicationError) Status() string {
	return e.ErrStatus
}

func NewError(message string, status string, err error) ApplicationError {
	return applicationError{
		ErrMessage: message,
		ErrStatus:  status,
		ErrError:   err,
	}
}

func NewErrorFromMessage(message string, status string) ApplicationError {
	return NewError(message, status, errors.New(message))
}

func NewErrorFromBytes(bytes []byte) (ApplicationError, error) {
	var appErr applicationError
	if err := json.Unmarshal(bytes, &appErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return appErr, nil
}
