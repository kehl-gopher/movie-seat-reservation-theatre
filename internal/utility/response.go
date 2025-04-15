package utility

type Response struct {
	StatusCode int         `json:"status_code"`
	Status     string      `json:"status"`
	Message    string      `json:"message,omitempty"`
	Error      string      `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

// I don't know what to do with this for now actually... it's probably just BS to sheeeeeeeeeeeeeesh
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
		message = "Internal Server Error"
		status = "error"
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
