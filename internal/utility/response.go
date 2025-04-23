package utility

import "net/http"

type Response struct {
	StatusCode int         `json:"status_code,omitempty"`
	Status     string      `json:"status,omitempty"`
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

func BuildErrorResponse(statusCode int, error error, message, status string) Response {
	return buildResponse(statusCode, status, message, error, nil, nil)
}

func BuildSuccessResponse(statusCode int, message string, data interface{}, pagination interface{}) Response {
	return buildResponse(statusCode, "success", message, nil, data, pagination)
}

func buildResponse(statusCode int, status, message string, err error, data interface{}, pagination interface{}) Response {
	var errMsg string

	if err != nil {
		errMsg = err.Error()
	}
	if statusCode == 500 {
		message = "error"
		errMsg = "Internal Server Error"
	}

	return Response{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
		Error:      errMsg,
	}
}

type ValidationError struct {
	Response
	Errors map[string]string `json:"errors"`
}

func NewValidationError() *ValidationError {
	return &ValidationError{}
}

func ValidationErrorResponse(errors map[string]string, v *ValidationError) *ValidationError {
	v.Response.StatusCode = http.StatusUnprocessableEntity
	v.Response.Status = http.StatusText(http.StatusUnprocessableEntity)
	v.Response.Message = "validation error"
	return v
}

func UnAuthorizedResponse(message string, status string) Response {
	return buildResponse(http.StatusUnauthorized, status, message, nil, nil, nil)
}
