package auroradnsclient

import "fmt"

// ErrorResponse describes the format of a generic AuroraDNS API error
type ErrorResponse struct {
	ErrorCode string `json:"error"`
	Message   string `json:"errormsg"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("%s - %s", e.ErrorCode, e.Message)
}
