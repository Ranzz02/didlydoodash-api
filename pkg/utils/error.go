package utils

type APIErrorStatus string

const (
	Error   APIErrorStatus = "error"
	Warning APIErrorStatus = "warning"
	Success APIErrorStatus = "success"
)

type APIError struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Status  APIErrorStatus `json:"status" default:"success"`
	Err     error          `json:"-"`
}

func (e APIError) Error() string {
	if e.Err == nil {
		return e.Message
	}
	return e.Err.Error()
}

func NewError(code int, message string, err error) APIError {
	return APIError{Code: code, Message: message, Status: Error, Err: err}
}

func NewWarning(code int, message string) APIError {
	return APIError{Code: code, Message: message, Status: Warning}
}

func NewSuccess(message string) APIError {
	return APIError{Code: 200, Message: message, Status: Success}
}
