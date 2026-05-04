package httpx

import "net/http"

type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func InvalidRequest(msg string, err error) *AppError {
	return &AppError{
		Status:  http.StatusBadRequest,
		Code:    "InvalidRequest",
		Message: msg,
		Err:     err,
	}
}

func NotFound(msg string, err error) *AppError {
	return &AppError{
		Status:  http.StatusNotFound,
		Code:    "NotFound",
		Message: msg,
		Err:     err,
	}
}

func ExternalAPI(msg string, err error) *AppError {
	return &AppError{
		Status:  http.StatusBadGateway,
		Code:    "ExternalAPIError",
		Message: msg,
		Err:     err,
	}
}

func Internal(msg string, err error) *AppError {
	return &AppError{
		Status:  http.StatusInternalServerError,
		Code:    "InternalError",
		Message: msg,
		Err:     err,
	}
}
