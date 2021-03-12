package api

type (
	ResponseError struct {
		Code    string            `json:"code"`
		Message string            `json:"message"`
		Reasons map[string]string `json:"reasons"`
		Details []interface{}     `json:"details,omitempty"`
	}
	Response struct {
		Status    int                    `json:"status"`
		Content   map[string]interface{} `json:"content,omitempty"`
		Error     *ResponseError         `json:"error,omitempty"`
	}
)

func NewResponse(status int, content map[string]interface{}, err *ResponseError) Response {
	if content == nil {
		content = make(map[string]interface{}, 0)
	}
	return Response{
		Status:  status,
		Content: content,
		Error:   err,
	}
}

func NewResponseError(code, message string, reasons map[string]string, details []interface{}) *ResponseError {
	if reasons == nil {
		reasons = make(map[string]string, 0)
	}
	if details == nil {
		details = make([]interface{}, 0)
	}
	return &ResponseError{
		Code:    code,
		Message: message,
		Reasons: reasons,
		Details: details,
	}
}