package secretshub

import "fmt"

func (e *ErrorResponse) String() string {
	if e == nil {
		return "No error response"
	}
	code := e.Code
	message := e.Message
	if e.Code == "" {
		code = "Empty error code"
	}
	if e.Message == "" {
		message = "Empty error message"
	}
	m := fmt.Sprintf("ErrorCode: %s, Message: %s", code, message)
	if e.Description != nil && *e.Description != "" {
		m += fmt.Sprintf(", Description: %s", *e.Description)
	}
	return m
}
