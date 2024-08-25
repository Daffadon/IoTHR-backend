package errors

import "fmt"

type ErrorInstance struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorInstance) Error() string {
	return fmt.Sprintf("Error code: %d, Error message: %s", e.Code, e.Message)
}
func (e *ErrorInstance) ReturnError(code int, message string) *ErrorInstance {
	return &ErrorInstance{
		Code:    code,
		Message: message,
	}
}
