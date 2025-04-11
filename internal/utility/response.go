package utility

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error,omitempty"`
	Message    string `json:"message"`
	Status     string `json:"status"`
}

type ValidationErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Error      []struct {
		Message string `json:"message"`
		Error   string `json:"error"`
	}
}

type SuccessResponse struct {
	StatusCode int                    `json:"status_code"`
	Status     string                 `json:"status"`
	Message    string                 `json:"message,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
	Pagination map[string]interface{} `json:"pagination,omitempty"`
}

func BuildValidationErrorResponse(statusCode int, err error, status string, in ...interface{}) []ErrorResponse {
	return nil
}
func BuildErrorResponse(statusCode int, error error, message, status string, in ...interface{}) ErrorResponse {
	return ErrorResponse{}
}
func BuildSuccessResponse(statusCode int, status, message string, data map[string]interface{}, pagination ...map[string]interface{}) SuccessResponse {
	return SuccessResponse{}
}
