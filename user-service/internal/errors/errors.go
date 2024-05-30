package errors

import "time"

const (
	TimeoutDuration = time.Second * 3
)

var ResponseRequestTimeout ErrorResponse = ErrorResponse{
	Message: "Request timeout",
}

type ErrorResponse struct {
	Message string `json:"message"`
}
